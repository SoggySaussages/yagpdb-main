package common

import "github.com/botlabs-gg/sgpdb/v2/lib/when/rules"

var All = []rules.Rule{
	SlashDMY(rules.Override),
}
