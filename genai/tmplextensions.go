package genai

import (
	"encoding/json"

	"github.com/botlabs-gg/yagpdb/v2/common/templates"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
)

func init() {
	templates.RegisterSetupFunc(func(ctx *templates.Context) {
		ctx.ContextFuncs["genaiComplete"] = tmplGenAIComplete(ctx)

		ctx.ContextFuncs["genaiComplete"] = tmplGenAIComplete(ctx)
		ctx.ContextFuncs["genaiCompleteComplex"] = tmplGenAICompleteComplex(ctx)
		ctx.ContextFuncs["genaiModerate"] = tmplGenAIModerate(ctx)
	})
}

func CreateGenAIInput(values ...interface{}) (*GenAIInput, error) {
	if len(values) < 1 {
		return &GenAIInput{}, nil
	}

	var m map[string]interface{}
	switch t := values[0].(type) {
	case templates.SDict:
		m = t
	case *templates.SDict:
		m = *t
	case map[string]interface{}:
		m = t
	case *GenAIInput:
		return t, nil
	default:
		dict, err := templates.StringKeyDictionary(values...)
		if err != nil {
			return nil, err
		}
		m = dict
	}

	encoded, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var input *GenAIInput
	err = json.Unmarshal(encoded, &input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func tmplGenAIComplete(ctx *templates.Context) interface{} {
	return func(systemMsg, userMsg string, maxTokens int64) (string, error) {
		if ctx.IncreaseCheckGenericAPICall() {
			return "", templates.ErrTooManyAPICalls
		}

		config, err := GetConfig(ctx.GS.ID)
		if err != nil {
			return "", err
		}

		provider := GenAIProviderFromID(config.Provider)
		response, _, err := provider.BasicCompletion(&dstate.GuildState{ID: ctx.GS.ID, OwnerID: ctx.GS.OwnerID}, systemMsg, userMsg, maxTokens, ctx.CurrentFrame.CS.NSFW)

		return response.Content, err
	}
}

func tmplGenAICompleteComplex(ctx *templates.Context) interface{} {
	return func(inputInterface interface{}) (*GenAIResponse, error) {
		if ctx.IncreaseCheckGenericAPICall() {
			return nil, templates.ErrTooManyAPICalls
		}

		input, err := CreateGenAIInput(inputInterface)
		if err != nil {
			return nil, err
		}

		config, err := GetConfig(ctx.GS.ID)
		if err != nil {
			return nil, err
		}

		provider := GenAIProviderFromID(config.Provider)
		nsfw := BotSystemMessagePromptAppendNonNSFW
		if ctx.CurrentFrame.CS.NSFW {
			nsfw = BotSystemMessagePromptAppendNSFW
		}
		input.BotSystemMessage = BotSystemMessagePromptGeneric + "\n" + nsfw
		resp, _, err := provider.ComplexCompletion(&dstate.GuildState{ID: ctx.GS.ID, OwnerID: ctx.GS.OwnerID}, input)
		return resp, err
	}
}

func tmplGenAIModerate(ctx *templates.Context) interface{} {
	return func(message string) (*GenAIModerationCategoryProbability, error) {
		if ctx.IncreaseCheckGenericAPICall() {
			return nil, templates.ErrTooManyAPICalls
		}

		config, err := GetConfig(ctx.GS.ID)
		if err != nil {
			return nil, err
		}

		provider := GenAIProviderFromID(config.Provider)
		response, _, err := provider.ModerateMessage(&dstate.GuildState{ID: ctx.GS.ID, OwnerID: ctx.GS.OwnerID}, message)

		return response, err
	}
}
