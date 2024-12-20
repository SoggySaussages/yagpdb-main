package genai

var DBSchemas = []string{`
CREATE TABLE IF NOT EXISTS genai_configs (
	guild_id BIGINT PRIMARY KEY,

	enabled BOOL NOT NULL,
	provider INT NOT NULL,
	model TEXT NOT NULL,
	key BYTEA NOT NULL,
	base_cmd_enabled BOOL NOT NULL
);
`, `
CREATE TABLE IF NOT EXISTS genai_commands (
	id BIGINT NOT NULL,
	guild_id BIGINT NOT NULL,

	enabled BOOL NOT NULL,

	triggers TEXT[] NOT NULL,
	prompt TEXT NOT NULL,
	allow_input BOOL NOT NULL,
	whitelisted_context BIGINT NOT NULL,
	max_tokens INT NOT NULL,

	autodelete_response BOOL NOT NULL,
	autodelete_trigger BOOL NOT NULL,

	autodelete_response_delay INT NOT NULL,
	autodelete_trigger_delay INT NOT NULL,

	channels BIGINT[],
	channels_whitelist_mode BOOL NOT NULL,

	roles BIGINT[],
	roles_whitelist_mode BOOL NOT NULL,

	PRIMARY KEY(guild_id, id)
);
`, `
CREATE INDEX IF NOT EXISTS genai_commands_guild_idx ON genai_commands(guild_id);
`}
