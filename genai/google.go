package genai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	google "cloud.google.com/go/vertexai/genai"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"google.golang.org/api/option"
)

type GenAIProviderGoogle struct{}

func (p GenAIProviderGoogle) ID() GenAIProviderID {
	return GenAIProviderGoogleID
}

func (p GenAIProviderGoogle) String() string {
	return "Google"
}

func (p GenAIProviderGoogle) DefaultModel() string {
	return "gemini-1.5-flash-002" // cheapest model as of Dec 2024
}

var GenAIModelMapGoogle = &GenAIProviderModelMap{
	"Gemini 1.0 Pro":   "gemini-1.0-pro-002",
	"Gemini 1.5 Pro":   "gemini-1.5-pro-002",
	"Gemini 1.5 Flash": "gemini-1.5-flash-002",
}

func (p GenAIProviderGoogle) ModelMap() *GenAIProviderModelMap {
	return GenAIModelMapGoogle
}

func (p GenAIProviderGoogle) KeyRequired() bool {
	return true
}

// ~ accurate as of Dec 2024
const CharacterCountToTokenRatioGoogle = 4 / 1

func (p GenAIProviderGoogle) CharacterTokenRatio() int {
	return CharacterCountToTokenRatioGoogle
}

type projectIDFromCredentials struct {
	ProjectID string `json:"project_id"`
}

func (p GenAIProviderGoogle) getCredentials(gs *dstate.GuildState) (projectID string, credentials []byte) {
	key, err := getAPIToken(gs)
	if err != nil {
		if err == ErrorNoAPIKey || err == ErrorAPIKeyInvalid {
			return "", nil
		}
		logger.Error(err)
		return "", nil
	}

	var projectIDStruct projectIDFromCredentials
	json.Unmarshal([]byte(key), &projectIDStruct)
	return projectIDStruct.ProjectID, []byte(key)
}

func (p GenAIProviderGoogle) EstimateTokens(combinedInput string, maxTokens int64) (inputEstimatedTokens, outputMaxTokens int64) {
	inputEstimatedTokens = int64(len(combinedInput) / CharacterCountToTokenRatioOpenAI)
	outputMaxTokens = maxTokens - inputEstimatedTokens
	return
}

func (p GenAIProviderGoogle) client(projectID string, credentials []byte) (*google.Client, error) {
	if projectID == "" || len(credentials) == 0 {
		return nil, errors.New("Your credentials are invalid.")
	}
	return google.NewClient(context.Background(), projectID, "us-central1", option.WithCredentialsJSON(credentials))
}

func (p GenAIProviderGoogle) ValidateAPIToken(gs *dstate.GuildState, token string) error {
	projectIDStruct := projectIDFromCredentials{}
	err := json.Unmarshal([]byte(token), &projectIDStruct)
	if err != nil {
		return fmt.Errorf("error unmarshalling to credentials file: %w", err)
	}

	client, err := p.client(projectIDStruct.ProjectID, []byte(token))
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(p.DefaultModel())
	gemini.SetMaxOutputTokens(1)
	_, err = gemini.GenerateContent(context.Background(), google.Text("1"))
	return err
}

func (p GenAIProviderGoogle) BasicCompletion(gs *dstate.GuildState, systemMsg, userMsg string, maxTokens int64, nsfw bool) (*GenAIResponse, *GenAIResponseUsage, error) {
	input := &GenAIInput{BotSystemMessage: BotSystemMessagePromptGeneric + BotSystemMessagePromptAppendSingleResponseContext, SystemMessage: systemMsg, UserMessage: userMsg, MaxTokens: maxTokens}
	if nsfw {
		input.BotSystemMessage = fmt.Sprintf("%s\n%s", input.BotSystemMessage, BotSystemMessagePromptAppendNSFW)
	} else {
		input.BotSystemMessage = fmt.Sprintf("%s\n%s", input.BotSystemMessage, BotSystemMessagePromptAppendNonNSFW)
	}
	return p.ComplexCompletion(gs, input)
}

func (p GenAIProviderGoogle) ComplexCompletion(gs *dstate.GuildState, input *GenAIInput) (*GenAIResponse, *GenAIResponseUsage, error) {
	config, err := GetConfig(gs.ID)
	if err != nil {
		return nil, nil, err
	}
	if !config.Enabled {
		return &GenAIResponse{}, &GenAIResponseUsage{}, err
	}

	client, err := p.client(p.getCredentials(gs))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(config.Model)
	presencePenalty := float32(0.05)
	gemini.GenerationConfig.PresencePenalty = &presencePenalty
	gemini.SetTemperature(1.1)
	gemini.SetMaxOutputTokens(int32(input.MaxTokens))
	gemini.SystemInstruction = &google.Content{Parts: []google.Part{google.Text(input.BotSystemMessage)}}

	if input.SystemMessage != "" {
		gemini.SystemInstruction.Parts = append(gemini.SystemInstruction.Parts, google.Text(input.SystemMessage))
	}

	prompt := google.Text("Please begin.")

	if input.UserMessage != "" {
		prompt = google.Text(input.UserMessage)
	}

	var tools []*google.FunctionDeclaration

	if input.Functions != nil {
		for _, fn := range *input.Functions {
			properties := make(map[string]*google.Schema, 0)
			for argName, argType := range fn.Arguments {
				var googleArgType google.Type
				switch argType {
				case "array":
					googleArgType = google.TypeArray
				case "bool", "boolean":
					googleArgType = google.TypeBoolean
				case "int", "integer":
					googleArgType = google.TypeInteger
				case "number":
					googleArgType = google.TypeNumber
				case "object":
					googleArgType = google.TypeObject
				case "string":
					googleArgType = google.TypeString
				default:
					googleArgType = google.TypeUnspecified
				}
				properties[argName] = &google.Schema{
					Type: googleArgType,
				}
			}

			tools = append(tools, &google.FunctionDeclaration{
				Name:        fn.Name,
				Description: fn.Definition,
				Parameters: &google.Schema{
					Type:       google.TypeObject,
					Properties: properties,
				}})
		}
	}

	if len(tools) > 0 {
		gemini.Tools = []*google.Tool{&google.Tool{
			FunctionDeclarations: tools,
		}}
	}

	resp, err := gemini.GenerateContent(context.Background(), prompt)
	if err != nil {
		logger.Error(err)
		return nil, nil, fmt.Errorf("error generating content: %w", err)
	}
	fmt.Println(json.Marshal(resp))

	choice := resp.Candidates[0]
	var content string
	for _, p := range choice.Content.Parts {
		c, ok := p.(google.Text)
		if ok {
			content += "\n" + string(c)
		}
	}

	var functionResponse []GenAIFunctionResponse
	if len(choice.FunctionCalls()) > 0 {
		for _, f := range choice.FunctionCalls() {
			currentFunc := GenAIFunctionResponse{}
			currentFunc.Name = f.Name
			currentFunc.Arguments = f.Args
			functionResponse = append(functionResponse, currentFunc)
		}
	}

	return &GenAIResponse{
			Content:   content,
			Functions: &functionResponse,
		}, &GenAIResponseUsage{
			InputTokens:  int64(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int64(resp.UsageMetadata.CandidatesTokenCount),
		}, nil
}

func (p GenAIProviderGoogle) ModerateMessage(gs *dstate.GuildState, message string) (*GenAIModerationCategoryProbability, *GenAIResponseUsage, error) {
	//	key, err := getAPIToken(gs)
	//	if err != nil {
	//		if err == ErrorNoAPIKey || err == ErrorAPIKeyInvalid {
	//			return &GenAIModerationCategoryProbability{}, nil, nil
	//		}
	//		return nil, nil, err
	//	}

	config, err := GetConfig(gs.ID)
	if err != nil {
		return nil, nil, err
	}

	client, err := p.client(p.getCredentials(gs))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(config.Model)
	presencePenalty := float32(0.05)
	gemini.GenerationConfig.PresencePenalty = &presencePenalty
	gemini.SetTemperature(1.1)
	gemini.GenerationConfig = google.GenerationConfig{}
	gemini.SetMaxOutputTokens(1)
	gemini.SafetySettings = []*google.SafetySetting{
		// define a bunch of safety settings with low thresholds so we get the
		// probability data in the response
		{Category: google.HarmCategoryHateSpeech, Threshold: google.HarmBlockLowAndAbove},
		{Category: google.HarmCategoryDangerousContent, Threshold: google.HarmBlockLowAndAbove},
		{Category: google.HarmCategoryHarassment, Threshold: google.HarmBlockLowAndAbove},
		{Category: google.HarmCategorySexuallyExplicit, Threshold: google.HarmBlockLowAndAbove},
	}
	resp, err := gemini.GenerateContent(context.Background(), google.Text(message))
	if err != nil {
		logger.Error(err)
		return nil, nil, nil
	}

	if resp.PromptFeedback == nil {
		// no categories were high enough to be blocked by the "low" threshold
		return &GenAIModerationCategoryProbability{}, &GenAIResponseUsage{
			InputTokens:  int64(resp.UsageMetadata.PromptTokenCount),
			OutputTokens: int64(resp.UsageMetadata.CandidatesTokenCount)}, nil
	}

	response := GenAIModerationCategoryProbability{}
	for _, r := range resp.PromptFeedback.SafetyRatings {
		switch r.Category {
		case google.HarmCategoryHateSpeech:
			response["Hate"] = float64(r.ProbabilityScore)
		case google.HarmCategoryDangerousContent:
			response["Violence"] = float64(r.ProbabilityScore)
		case google.HarmCategoryHarassment:
			response["Harassment"] = float64(r.ProbabilityScore)
		case google.HarmCategorySexuallyExplicit:
			response["Sexual"] = float64(r.ProbabilityScore)
		}
	}

	return &response, &GenAIResponseUsage{
		InputTokens:  int64(resp.UsageMetadata.PromptTokenCount),
		OutputTokens: int64(resp.UsageMetadata.CandidatesTokenCount)}, nil
}

var GenAIProviderGoogleWebData = &GenAIProviderWebDescriptions{
	ObtainingAPIKeyInstructions: template.HTML(`Step one: Create an account.
	<br>
	Visit <a href="https://platform.openai.com/docs/guides/production-best-practices/api-keys#setting-up-your-organization">Google's website</a> to do this.
	<br>
	<br>
	Step two: Set up payment method.
	<br>
	You must set up a payment method in order to make requests to Google. Do so on <a href="https://platform.openai.com/settings/organization/billing/overview">Google's API dashboard</a>.
	<br>
	<br>
	Step three: Set a Budget Limit.
	<br>
	You must set a monthly budget limit within reason to prevent yourself from going into credit debt with Google. Do so on <a href="https://platform.openai.com/settings/organization/limits">Google's API dashboard</a>.
	<br>
	<br>
	Step four: Create an API key.
	<br>
	Create an API key on <a href="https://platform.openai.com/api-keys">Google's Dashboard</a>. Set the mode to <strong>restricted</strong>, set every permission to <strong>None</strong>, and then set the "Model capabilities" permission to <strong>Write</strong>.
	<br>
	<br>
	Step five: Copy the API key to YAGPDB.
	<br>
	Click copy, then paste the new API key into the "API Key" field on this page.`),
	ModelDescriptionsURL: "https://platform.openai.com/docs/models",
	ModelForModeration:   "omni-moderation-latest",
}

func (p GenAIProviderGoogle) WebData() *GenAIProviderWebDescriptions {
	return GenAIProviderGoogleWebData
}
