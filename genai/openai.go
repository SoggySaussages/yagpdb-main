package genai

import (
	"context"
	"encoding/json"

	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type GenAIProviderOpenAI struct{}

func (p GenAIProviderOpenAI) Type() GenAIProviderID {
	return GenAIProviderOpenAIID
}

func (p GenAIProviderOpenAI) DefaultModel() string {
	return openai.ChatModelGPT4oMini // cheapest model as of Dec 2024
}

var GenAIModelMapOpenAI = &GenAIProviderModelMap{
	"o1 Preview":    openai.ChatModelO1Preview,
	"o1 Mini":       openai.ChatModelO1Mini,
	"GPT 4o":        openai.ChatModelGPT4o,
	"ChatGPT 4o":    openai.ChatModelChatgpt4oLatest,
	"GPT 4o Mini":   openai.ChatModelGPT4oMini,
	"GPT 4 Turbo":   openai.ChatModelGPT4Turbo,
	"GPT 4":         openai.ChatModelGPT4,
	"GPT 3.5 Turbo": openai.ChatModelGPT3_5Turbo,
}

func (p *GenAIProviderOpenAI) ModelMap() *GenAIProviderModelMap {
	return GenAIModelMapOpenAI
}

func (p *GenAIProviderOpenAI) ComplexCompletion(gs dstate.GuildState, input *GenAIInput) (*GenAIResponse, *GenAIResponseUsage, error) {
	var messages []openai.ChatCompletionMessageParamUnion

	if input.SystemMessage != "" {
		messages = append(messages, openai.SystemMessage(input.SystemMessage))
	}

	if input.UserMessage != "" {
		messages = append(messages, openai.UserMessage(input.UserMessage))
	}

	var tools []openai.ChatCompletionToolParam

	if input.Functions != nil {
		for _, fn := range *input.Functions {
			properties := make(map[string]interface{}, 0)
			for argName, argType := range fn.Arguments {
				properties[argName] = map[string]string{
					"type": argType,
				}
			}

			tools = append(tools, openai.ChatCompletionToolParam{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name:        openai.String(fn.Name),
					Description: openai.String(fn.Definition),
					Parameters: openai.F(openai.FunctionParameters{
						"type":       "object",
						"properties": properties,
					}),
				}),
			})
		}
	}

	requestParams := openai.ChatCompletionNewParams{}

	if len(messages) > 0 {
		requestParams.Messages = openai.F(messages)
	}

	if len(tools) > 0 {
		requestParams.Tools = openai.F(tools)
	}

	key, err := getAPIToken(&gs)
	if err != nil {
		return nil, nil, err
	}

	client := openai.NewClient(option.WithAPIKey(key))

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
			openai.UserMessage("Say this is a test"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return nil, nil, err
	}

	choice := chatCompletion.Choices[0]
	content := choice.Message.Content
	if choice.Message.Refusal != "" {
		content = choice.Message.Refusal
	}

	var functionResponse []GenAIFunctionResponse
	if len(choice.Message.ToolCalls) > 0 {
		currentFunc := GenAIFunctionResponse{}
		functionCall := choice.Message.ToolCalls[0].Function
		currentFunc.Name = functionCall.Name
		json.Unmarshal([]byte(functionCall.Arguments), &currentFunc.Arguments)
		functionResponse = append(functionResponse, currentFunc)
	}

	return &GenAIResponse{
			Content:   content,
			Functions: &functionResponse,
		}, &GenAIResponseUsage{
			InputTokens:  uint64(chatCompletion.Usage.PromptTokens),
			OutputTokens: uint64(chatCompletion.Usage.CompletionTokens),
		}, nil
}
