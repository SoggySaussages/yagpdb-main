package util

import (
	"github.com/SoggySaussages/syzygy/bot"
	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/common"
	"github.com/SoggySaussages/syzygy/lib/dcmd"
)

func isExecedByCC(data *dcmd.Data) bool {
	if v := data.Context().Value(commands.CtxKeyExecutedByCC); v != nil {
		if cast, _ := v.(bool); cast {
			return true
		}
	}

	return false
}

func RequireBotAdmin(inner dcmd.RunFunc) dcmd.RunFunc {
	return func(data *dcmd.Data) (interface{}, error) {
		//	if isExecedByCC(data) {
		//		return "", nil
		//	}

		if admin, err := bot.IsBotAdmin(data.Author.ID); admin && err == nil {
			return inner(data)
		}

		return "", nil
	}
}

func RequireOwner(inner dcmd.RunFunc) dcmd.RunFunc {
	return func(data *dcmd.Data) (interface{}, error) {
		//	if isExecedByCC(data) {
		//		return " ", nil
		//	}

		if common.IsOwner(data.Author.ID) {
			return inner(data)
		}

		return "", nil
	}
}
