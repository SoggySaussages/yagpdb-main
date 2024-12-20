package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
)

type GenAIProviderAnthropic struct{}

func (p GenAIProviderAnthropic) ID() GenAIProviderID {
	return GenAIProviderAnthropicID
}

func (p GenAIProviderAnthropic) String() string {
	return "Anthropic"
}

func (p GenAIProviderAnthropic) DefaultModel() string {
	return anthropic.ModelClaude_3_Haiku_20240307 // cheapest model as of Dec 2024
}

var GenAIModelMapAnthropic = &GenAIProviderModelMap{
	"Claude 3 Haiku":    anthropic.ModelClaude_3_Haiku_20240307,
	"Claude 3 Sonnet":   anthropic.ModelClaude_3_Sonnet_20240229,
	"Claude 3 Opus":     anthropic.ModelClaude3OpusLatest,
	"Claude 3.5 Haiku":  anthropic.ModelClaude3_5HaikuLatest,
	"Claude 3.5 Sonnet": anthropic.ModelClaude3_5SonnetLatest,
}

func (p GenAIProviderAnthropic) ModelMap() *GenAIProviderModelMap {
	return GenAIModelMapAnthropic
}

func (p GenAIProviderAnthropic) KeyRequired() bool {
	return true
}

// ~ accurate for English text as of Dec 2024
const CharacterCountToTokenRatioAnthropic = 4 / 1

func (p GenAIProviderAnthropic) CharacterTokenRatio() int {
	return CharacterCountToTokenRatioAnthropic
}

func (p GenAIProviderAnthropic) EstimateTokens(combinedInput string, maxTokens int64) (inputEstimatedTokens, outputMaxTokens int64) {
	inputEstimatedTokens = int64(len(combinedInput) / CharacterCountToTokenRatioAnthropic)
	outputMaxTokens = maxTokens - inputEstimatedTokens
	return
}

func (p GenAIProviderAnthropic) ValidateAPIToken(gs *dstate.GuildState, token string) error {
	// make a really cheap call to test the key
	client := anthropic.NewClient(option.WithAPIKey(token))
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Messages:  anthropic.F([]anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("1"))}),
		Model:     anthropic.F(p.DefaultModel()),
		MaxTokens: anthropic.Int(1),
	})
	return err
}

func (p GenAIProviderAnthropic) BasicCompletion(gs *dstate.GuildState, systemMsg, userMsg string, maxTokens int64, nsfw bool) (*GenAIResponse, *GenAIResponseUsage, error) {
	input := &GenAIInput{BotSystemMessage: BotSystemMessagePromptGeneric + BotSystemMessagePromptAppendSingleResponseContext, SystemMessage: systemMsg, UserMessage: userMsg, MaxTokens: maxTokens}
	if nsfw {
		input.BotSystemMessage = fmt.Sprintf("%s\n%s", input.BotSystemMessage, BotSystemMessagePromptAppendNSFW)
	} else {
		input.BotSystemMessage = fmt.Sprintf("%s\n%s", input.BotSystemMessage, BotSystemMessagePromptAppendNonNSFW)
	}
	return p.ComplexCompletion(gs, input)
}

func (p GenAIProviderAnthropic) convertToJSONSchema(args json.RawMessage) interface{} {
	return json.RawMessage(fmt.Sprintf(`{"$schema": "http://json-schema.org/draft/2020-12/schema",
	  "properties": %s,
	  "type": "object",
	  "additional_properties": false,
	  "required": []
	}`, string(args)))
}

func (p GenAIProviderAnthropic) ComplexCompletion(gs *dstate.GuildState, input *GenAIInput) (*GenAIResponse, *GenAIResponseUsage, error) {
	key, err := getAPIToken(gs)
	if err != nil {
		if err == ErrorNoAPIKey {
			return &GenAIResponse{Content: "Please set your API key on the dashboard to use Generative AI."}, &GenAIResponseUsage{}, nil
		}
		if err == ErrorAPIKeyInvalid {
			return &GenAIResponse{Content: err.Error()}, &GenAIResponseUsage{}, nil
		}
		return nil, nil, err
	}

	config, err := GetConfig(gs.ID)
	if err != nil {
		return nil, nil, err
	}

	systemMessages := []anthropic.TextBlockParam{}

	systemMessages = append(systemMessages, anthropic.NewTextBlock(input.BotSystemMessage))

	if input.SystemMessage != "" {
		systemMessages = append(systemMessages, anthropic.NewTextBlock(input.SystemMessage))
	}

	var tools []anthropic.ToolParam

	if input.Functions != nil {
		for _, fn := range *input.Functions {
			properties := make(map[string]interface{}, 0)
			for argName, argType := range fn.Arguments {
				properties[argName] = map[string]string{
					"type": argType,
				}
			}

			inputSchema, _ := json.Marshal(properties)
			inSch := p.convertToJSONSchema(inputSchema)
			tools = append(tools, anthropic.ToolParam{
				Name:        anthropic.String(fn.Name),
				Description: anthropic.String(fn.Description),
				InputSchema: anthropic.F(inSch),
			})
		}
	}

	model := config.Model
	if input.ModelOverride != "" {
		for _, v := range *p.ModelMap() {
			if v == input.ModelOverride {
				model = input.ModelOverride
				break
			}
		}
	}

	requestParams := anthropic.MessageNewParams{Model: anthropic.F(model), MaxTokens: anthropic.Int(input.MaxTokens), System: anthropic.F(systemMessages), Messages: anthropic.F([]anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("Please begin."))}), Temperature: anthropic.Float(1)}

	if input.UserMessage != "" {
		requestParams.Messages = anthropic.F([]anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(input.UserMessage))})
	}

	if len(tools) > 0 {
		requestParams.Tools = anthropic.F(tools)
	}

	client := anthropic.NewClient(option.WithAPIKey(key))

	messageResp, err := client.Messages.New(context.Background(), requestParams)
	if err != nil {
		return nil, nil, err
	}

	content := messageResp.Content
	var functionResponse []GenAIFunctionResponse
	var contentString string

	for _, block := range content {
		if block.Type == anthropic.ContentBlockTypeToolUse {
			currentFunc := GenAIFunctionResponse{
				Name: block.Name,
			}
			json.Unmarshal(block.Input, &currentFunc.Arguments)
			functionResponse = append(functionResponse, currentFunc)
		} else {
			contentString = block.Text
		}
	}

	return &GenAIResponse{
			Content:   contentString,
			Functions: &functionResponse,
		}, &GenAIResponseUsage{
			InputTokens:  int64(messageResp.Usage.InputTokens),
			OutputTokens: int64(messageResp.Usage.OutputTokens),
		}, nil
}

func (p GenAIProviderAnthropic) ModerateMessage(gs *dstate.GuildState, message string) (*GenAIModerationCategoryProbability, *GenAIResponseUsage, error) {
	input := CustomModerateFunction
	input.UserMessage = message
	input.MaxTokens = 96
	input.ModelOverride = anthropic.ModelClaude3_5HaikuLatest

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

var GenAIProviderAnthropicWebData = &GenAIProviderWebDescriptions{
	ObtainingAPIKeyInstructions: template.HTML(`Step one: Create an account.
	<br>
	Visit <a href="https://console.anthropic.com">Anthropic's website</a> to do this. Once you've created your account, you'll be prompted to give your name and organization name.
	<br>
	<br>
	Step two: Set up payment method.
	<br>
	You must set up a payment method in order to make requests to Anthropic. Do so on <a href="https://console.anthropic.com/settings/billing">Anthropic's API dashboard</a>. You will be prompted to provide detailed information on your organization.
	<br>
	<br>
	Step three: Set a Budget Limit.
	<br>
	You must set a monthly budget limit within reason to prevent yourself from going into credit debt with Anthropic. Do so on <a href="https://console.anthropic.com/settings/limits">Anthropic's API dashboard</a>, scroll to the bottom under <strong>Monthly limit</strong> and click <strong>Change limit</strong>.
	<br>
	<br>
	Step four: Create an API key.
	<br>
	Create an API key on <a href="https://console.anthropic.com/settings/keys">Anthropic's Dashboard</a>. Give it a name and click <strong>Create</strong>.
	<br>
	<br>
	Step five: Copy the API key to YAGPDB.
	<br>
	Click copy, then paste the new API key into the "API Key" field on this page.`),
	ModelDescriptionsURL: "https://platform.anthropic.com/docs/models",
	ModelForModeration:   anthropic.ModelClaude3_5HaikuLatest,
	PlaygroundURL:        "https://console.anthropic.com/workbench",
}

func (p GenAIProviderAnthropic) WebData() *GenAIProviderWebDescriptions {
	return GenAIProviderAnthropicWebData
}
