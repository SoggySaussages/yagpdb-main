package verification

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/common/config"
)

var confGoogleReCAPTCHASiteKey = config.RegisterOption("sgpdb.google.recaptcha_site_key", "Google reCAPTCHA site key", "")
var confGoogleReCAPTCHASecret = config.RegisterOption("sgpdb.google.recaptcha_secret", "Google reCAPTCHA site secret", "")
var confVerificationTrackIPs = config.RegisterOption("sgpdb.verification.track_ips", "Track verified users ip", true)

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Verification",
		SysName:  "verification",
		Category: common.PluginCategoryMisc,
	}
}

var logger = common.GetPluginLogger(&Plugin{})

func RegisterPlugin() {

	if confGoogleReCAPTCHASecret.GetString() == "" || confGoogleReCAPTCHASiteKey.GetString() == "" {
		logger.Warn("no SGPDB_GOOGLE_RECAPTCHA_SECRET and/or SGPDB_GOOGLE_RECAPTCHA_SITE_KEY provided, not enabling verification plugin")
		return
	}

	common.InitSchemas("verification", DBSchemas...)

	common.RegisterPlugin(&Plugin{})
}

const (
	DefaultPageContent = `## Verification

Please solve the following reCAPTCHA to make sure you're not a robot`
)

const DefaultDMMessage = `{{sendMessage nil (cembed
"title" "Are you a bot?"
"description" (printf "Please solve the CAPTCHA at this link to make sure you're human, before you can enter %s: %s" .Server.Name .Link)
)}}`
