package genai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strconv"

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

	model := config.Model
	if input.ModelOverride != "" {
		for _, v := range *p.ModelMap() {
			if v == input.ModelOverride {
				model = input.ModelOverride
				break
			}
		}
	}

	gemini := client.GenerativeModel(model)
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
				Description: fn.Description,
				Parameters: &google.Schema{
					Type:       google.TypeObject,
					Properties: properties,
				}})
		}
	}

	if len(tools) > 0 {
		gemini.Tools = []*google.Tool{
			{FunctionDeclarations: tools},
		}
	}

	resp, err := gemini.GenerateContent(context.Background(), prompt)
	if err != nil {
		logger.Error(err)
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
	input := CustomModerateFunction
	input.UserMessage = message
	input.MaxTokens = 96

	r, u, err := p.ComplexCompletion(gs, &input)
	if err != nil {
		logger.Error(err)
		return &GenAIModerationCategoryProbability{}, u, nil
	}

	if len(*r.Functions) == 0 {
		return &GenAIModerationCategoryProbability{}, u, nil
	}

	modResp := (*r.Functions)[0]
	if len(modResp.Arguments) == 0 {
		return &GenAIModerationCategoryProbability{}, u, nil
	}

	response := GenAIModerationCategoryProbability{}
	for cat, prob := range modResp.Arguments {
		probInt := 0
		t := reflect.ValueOf(prob)
		switch {
		case t.CanInt():
			probInt = int(t.Int())
		case t.CanFloat():
			probInt = int(t.Float())
		case t.CanUint():
			probInt = int(t.Uint())
		case t.Kind() == reflect.String:
			parsed, _ := strconv.ParseInt(t.String(), 10, 64)
			probInt = int(parsed)
		}

		response[cat] = float64(probInt) / 100.0
	}

	return &response, u, nil
}

var GenAIProviderGoogleWebData = &GenAIProviderWebDescriptions{
	ObtainingAPIKeyInstructions: template.HTML(`Step one: Create an account.
	<br>
	Visit <a href="https://console.cloud.google.com">Google's website</a> to do this.
	<br>
	<br>
	Step two: Set up project.
	<br>
	On the same page, click <strong>Select a project</strong> in the top left corner of the screen, then click <strong>New Project</strong>. Finally, name your project (optional) and click <strong>Create</strong>.
	<br>
	<br>
	Step three: Set up payment method.
	<br>
	Go to <a href="https://console.cloud.google.com/billing">your account's billing page</a>, click <strong>Add Billing Account</strong>, and complete the process.
	<br>
	Once done, go to <a href="https://console.cloud.google.com/billing/projects">your projects billing dasboard</a> then click the three dots ... on your project, and click <strong>change billing</strong>. Select the billing account you just created from the dropdown, and click <strong>Set Account</strong>.
	<br>
	<br>
	Step four: Set a Budget Limit.
	<br>
	You must set a monthly budget alert within reason to protect yourself from going into credit debt with Google. Do so on <a href="https://console.cloud.google.com/billing">Google's Billing dashboard</a>.
	<br>
	<br>
	Step five: Enable the AI API.
	<br>
	Enable the Vertex AI API plugin to be able to make requests to it. Do so on <a href="https://console.cloud.google.com/apis/library/aiplatform.googleapis.com?project=our-card-445307-d6">Google's API Library</a>.
	<br>
	<br>
	Step six: Create a service account.
	<br>
	Create a service account to allow the bot to make requests to Vertex AI API. Do so on <a href="https://console.cloud.google.com/iam-admin/serviceaccounts/create">the Google Cloud Dashboard</a>.
	<br>
	<br>
	Step seven: Create a credentials file.
	Once you have created your service account, you'll be redirected to the <a href="https://console.cloud.google.com/iam-admin/serviceaccounts">service accounts page</a>. From there, click the three dots ... on your service account, and click <strong>Manage Keys</strong>.
	<br>
	Click <strong>Add Key</strong>, then <strong>Create New Key</strong>, select <strong>JSON</strong>, and then <strong>Create</strong>.
	<br>
	<br>
	Step eight: Open your credentials file.
	<br>
	After creating your new key, the credentials file will be downloaded to your device. Find that file and open it. You should see a bunch of text starting with <code>"type": "service_account"</code>.
	<br>
	<br>
	Step nine: Copy the API key to YAGPDB.
	<br>
	You must select the <strong>entire contents</strong> of your file (not the name, all the text inside the file), and then copy it. Then, paste the <strong>entire</strong> file into the "API Key" field on this page.`),
	ModelDescriptionsURL: "https://ai.google.dev/pricing",
	ModelForModeration:   "(whichever model you choose)",
	PlaygroundURL:        "https://aistudio.google.com/prompts/new_chat",
}

func (p GenAIProviderGoogle) WebData() *GenAIProviderWebDescriptions {
	return GenAIProviderGoogleWebData
}
