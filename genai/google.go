package genai

import (
	"context"
	"fmt"
	"html/template"

	google "cloud.google.com/go/vertexai/genai"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type GenAIProviderGoogle struct{}

func (p GenAIProviderGoogle) ID() GenAIProviderID {
	return GenAIProviderGoogleID
}

func (p GenAIProviderGoogle) String() string {
	return "Google"
}

func (p GenAIProviderGoogle) DefaultModel() string {
	return openai.ChatModelGPT4oMini // cheapest model as of Dec 2024
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

// ~ accurate for English text as of Dec 2024
const CharacterCountToTokenRatioGoogle = 4 / 1

func (p GenAIProviderGoogle) CharacterTokenRatio() int {
	return CharacterCountToTokenRatioGoogle
}

func (p GenAIProviderGoogle) EstimateTokens(combinedInput string, maxTokens int64) (inputEstimatedTokens, outputMaxTokens int64) {
	inputEstimatedTokens = int64(len(combinedInput) / CharacterCountToTokenRatioGoogle)
	outputMaxTokens = maxTokens - inputEstimatedTokens
	return
}

func (p GenAIProviderGoogle) ValidateAPIToken(gs *dstate.GuildState, token string) error {
	// make a really cheap (%0.02 of a cent) call to test the key
	client := openai.NewClient(option.WithAPIKey(token))
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages:            openai.F([]openai.ChatCompletionMessageParamUnion{openai.UserMessage("1")}),
		Model:               openai.F(p.DefaultModel()),
		MaxCompletionTokens: openai.Int(1),
	})
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
	//	key, err := getAPIToken(gs)
	//	if err != nil {
	//		if err == ErrorNoAPIKey {
	//			return &GenAIResponse{Content: "Please set your API key on the dashboard to use Generative AI."}, &GenAIResponseUsage{}, nil
	//		}
	//		if err == ErrorAPIKeyInvalid {
	//			return &GenAIResponse{Content: err.Error()}, &GenAIResponseUsage{}, nil
	//		}
	//		return nil, nil, err
	//	}

	config, err := GetConfig(gs.ID)
	if err != nil {
		return nil, nil, err
	}

	client, err := google.NewClient(context.Background(), "yagpdb-394902", "us-central1")
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
		return nil, nil, fmt.Errorf("error generating content: %w", err)
	}

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

	client, err := google.NewClient(context.Background(), "yagpdb-394902", "us-central1")
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
	resp, err := gemini.GenerateContent(context.Background(), google.Text(message))
	if err != nil {
		return nil, nil, err
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
