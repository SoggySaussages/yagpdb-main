package common

import "github.com/SoggySaussages/sgpdb/lib/when/rules"

var All = []rules.Rule{
	SlashDMY(rules.Override),
}
