// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
	"github.com/volatiletech/sqlboiler/types"
)

// GenaiCommand is an object representing the database table.
type GenaiCommand struct {
	ID                      int64             `boil:"id" json:"id" toml:"id" yaml:"id"`
	GuildID                 int64             `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	Enabled                 bool              `boil:"enabled" json:"enabled" toml:"enabled" yaml:"enabled"`
	Triggers                types.StringArray `boil:"triggers" json:"triggers" toml:"triggers" yaml:"triggers"`
	Prompt                  string            `boil:"prompt" json:"prompt" toml:"prompt" yaml:"prompt"`
	AllowInput              bool              `boil:"allow_input" json:"allow_input" toml:"allow_input" yaml:"allow_input"`
	WhitelistedContext      int64             `boil:"whitelisted_context" json:"whitelisted_context" toml:"whitelisted_context" yaml:"whitelisted_context"`
	MaxTokens               int               `boil:"max_tokens" json:"max_tokens" toml:"max_tokens" yaml:"max_tokens"`
	AutodeleteResponse      bool              `boil:"autodelete_response" json:"autodelete_response" toml:"autodelete_response" yaml:"autodelete_response"`
	AutodeleteTrigger       bool              `boil:"autodelete_trigger" json:"autodelete_trigger" toml:"autodelete_trigger" yaml:"autodelete_trigger"`
	AutodeleteResponseDelay int               `boil:"autodelete_response_delay" json:"autodelete_response_delay" toml:"autodelete_response_delay" yaml:"autodelete_response_delay"`
	AutodeleteTriggerDelay  int               `boil:"autodelete_trigger_delay" json:"autodelete_trigger_delay" toml:"autodelete_trigger_delay" yaml:"autodelete_trigger_delay"`
	Channels                types.Int64Array  `boil:"channels" json:"channels,omitempty" toml:"channels" yaml:"channels,omitempty"`
	ChannelsWhitelistMode   bool              `boil:"channels_whitelist_mode" json:"channels_whitelist_mode" toml:"channels_whitelist_mode" yaml:"channels_whitelist_mode"`
	Roles                   types.Int64Array  `boil:"roles" json:"roles,omitempty" toml:"roles" yaml:"roles,omitempty"`
	RolesWhitelistMode      bool              `boil:"roles_whitelist_mode" json:"roles_whitelist_mode" toml:"roles_whitelist_mode" yaml:"roles_whitelist_mode"`

	R *genaiCommandR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L genaiCommandL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var GenaiCommandColumns = struct {
	ID                      string
	GuildID                 string
	Enabled                 string
	Triggers                string
	Prompt                  string
	AllowInput              string
	WhitelistedContext      string
	MaxTokens               string
	AutodeleteResponse      string
	AutodeleteTrigger       string
	AutodeleteResponseDelay string
	AutodeleteTriggerDelay  string
	Channels                string
	ChannelsWhitelistMode   string
	Roles                   string
	RolesWhitelistMode      string
}{
	ID:                      "id",
	GuildID:                 "guild_id",
	Enabled:                 "enabled",
	Triggers:                "triggers",
	Prompt:                  "prompt",
	AllowInput:              "allow_input",
	WhitelistedContext:      "whitelisted_context",
	MaxTokens:               "max_tokens",
	AutodeleteResponse:      "autodelete_response",
	AutodeleteTrigger:       "autodelete_trigger",
	AutodeleteResponseDelay: "autodelete_response_delay",
	AutodeleteTriggerDelay:  "autodelete_trigger_delay",
	Channels:                "channels",
	ChannelsWhitelistMode:   "channels_whitelist_mode",
	Roles:                   "roles",
	RolesWhitelistMode:      "roles_whitelist_mode",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpertypes_StringArray struct{ field string }

func (w whereHelpertypes_StringArray) EQ(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertypes_StringArray) NEQ(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertypes_StringArray) LT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_StringArray) LTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_StringArray) GT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_StringArray) GTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}

type whereHelpertypes_Int64Array struct{ field string }

func (w whereHelpertypes_Int64Array) EQ(x types.Int64Array) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_Int64Array) NEQ(x types.Int64Array) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_Int64Array) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_Int64Array) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpertypes_Int64Array) LT(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_Int64Array) LTE(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_Int64Array) GT(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_Int64Array) GTE(x types.Int64Array) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var GenaiCommandWhere = struct {
	ID                      whereHelperint64
	GuildID                 whereHelperint64
	Enabled                 whereHelperbool
	Triggers                whereHelpertypes_StringArray
	Prompt                  whereHelperstring
	AllowInput              whereHelperbool
	WhitelistedContext      whereHelperint64
	MaxTokens               whereHelperint
	AutodeleteResponse      whereHelperbool
	AutodeleteTrigger       whereHelperbool
	AutodeleteResponseDelay whereHelperint
	AutodeleteTriggerDelay  whereHelperint
	Channels                whereHelpertypes_Int64Array
	ChannelsWhitelistMode   whereHelperbool
	Roles                   whereHelpertypes_Int64Array
	RolesWhitelistMode      whereHelperbool
}{
	ID:                      whereHelperint64{field: "\"genai_commands\".\"id\""},
	GuildID:                 whereHelperint64{field: "\"genai_commands\".\"guild_id\""},
	Enabled:                 whereHelperbool{field: "\"genai_commands\".\"enabled\""},
	Triggers:                whereHelpertypes_StringArray{field: "\"genai_commands\".\"triggers\""},
	Prompt:                  whereHelperstring{field: "\"genai_commands\".\"prompt\""},
	AllowInput:              whereHelperbool{field: "\"genai_commands\".\"allow_input\""},
	WhitelistedContext:      whereHelperint64{field: "\"genai_commands\".\"whitelisted_context\""},
	MaxTokens:               whereHelperint{field: "\"genai_commands\".\"max_tokens\""},
	AutodeleteResponse:      whereHelperbool{field: "\"genai_commands\".\"autodelete_response\""},
	AutodeleteTrigger:       whereHelperbool{field: "\"genai_commands\".\"autodelete_trigger\""},
	AutodeleteResponseDelay: whereHelperint{field: "\"genai_commands\".\"autodelete_response_delay\""},
	AutodeleteTriggerDelay:  whereHelperint{field: "\"genai_commands\".\"autodelete_trigger_delay\""},
	Channels:                whereHelpertypes_Int64Array{field: "\"genai_commands\".\"channels\""},
	ChannelsWhitelistMode:   whereHelperbool{field: "\"genai_commands\".\"channels_whitelist_mode\""},
	Roles:                   whereHelpertypes_Int64Array{field: "\"genai_commands\".\"roles\""},
	RolesWhitelistMode:      whereHelperbool{field: "\"genai_commands\".\"roles_whitelist_mode\""},
}

// GenaiCommandRels is where relationship names are stored.
var GenaiCommandRels = struct {
}{}

// genaiCommandR is where relationships are stored.
type genaiCommandR struct {
}

// NewStruct creates a new relationship struct
func (*genaiCommandR) NewStruct() *genaiCommandR {
	return &genaiCommandR{}
}

// genaiCommandL is where Load methods for each relationship are stored.
type genaiCommandL struct{}

var (
	genaiCommandAllColumns            = []string{"id", "guild_id", "enabled", "triggers", "prompt", "allow_input", "whitelisted_context", "max_tokens", "autodelete_response", "autodelete_trigger", "autodelete_response_delay", "autodelete_trigger_delay", "channels", "channels_whitelist_mode", "roles", "roles_whitelist_mode"}
	genaiCommandColumnsWithoutDefault = []string{"id", "guild_id", "enabled", "triggers", "prompt", "allow_input", "whitelisted_context", "max_tokens", "autodelete_response", "autodelete_trigger", "autodelete_response_delay", "autodelete_trigger_delay", "channels", "channels_whitelist_mode", "roles", "roles_whitelist_mode"}
	genaiCommandColumnsWithDefault    = []string{}
	genaiCommandPrimaryKeyColumns     = []string{"guild_id", "id"}
)

type (
	// GenaiCommandSlice is an alias for a slice of pointers to GenaiCommand.
	// This should generally be used opposed to []GenaiCommand.
	GenaiCommandSlice []*GenaiCommand

	genaiCommandQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	genaiCommandType                 = reflect.TypeOf(&GenaiCommand{})
	genaiCommandMapping              = queries.MakeStructMapping(genaiCommandType)
	genaiCommandPrimaryKeyMapping, _ = queries.BindMapping(genaiCommandType, genaiCommandMapping, genaiCommandPrimaryKeyColumns)
	genaiCommandInsertCacheMut       sync.RWMutex
	genaiCommandInsertCache          = make(map[string]insertCache)
	genaiCommandUpdateCacheMut       sync.RWMutex
	genaiCommandUpdateCache          = make(map[string]updateCache)
	genaiCommandUpsertCacheMut       sync.RWMutex
	genaiCommandUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single genaiCommand record from the query using the global executor.
func (q genaiCommandQuery) OneG(ctx context.Context) (*GenaiCommand, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single genaiCommand record from the query.
func (q genaiCommandQuery) One(ctx context.Context, exec boil.ContextExecutor) (*GenaiCommand, error) {
	o := &GenaiCommand{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for genai_commands")
	}

	return o, nil
}

// AllG returns all GenaiCommand records from the query using the global executor.
func (q genaiCommandQuery) AllG(ctx context.Context) (GenaiCommandSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all GenaiCommand records from the query.
func (q genaiCommandQuery) All(ctx context.Context, exec boil.ContextExecutor) (GenaiCommandSlice, error) {
	var o []*GenaiCommand

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to GenaiCommand slice")
	}

	return o, nil
}

// CountG returns the count of all GenaiCommand records in the query, and panics on error.
func (q genaiCommandQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all GenaiCommand records in the query.
func (q genaiCommandQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count genai_commands rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q genaiCommandQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q genaiCommandQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if genai_commands exists")
	}

	return count > 0, nil
}

// GenaiCommands retrieves all the records using an executor.
func GenaiCommands(mods ...qm.QueryMod) genaiCommandQuery {
	mods = append(mods, qm.From("\"genai_commands\""))
	return genaiCommandQuery{NewQuery(mods...)}
}

// FindGenaiCommandG retrieves a single record by ID.
func FindGenaiCommandG(ctx context.Context, guildID int64, iD int64, selectCols ...string) (*GenaiCommand, error) {
	return FindGenaiCommand(ctx, boil.GetContextDB(), guildID, iD, selectCols...)
}

// FindGenaiCommand retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGenaiCommand(ctx context.Context, exec boil.ContextExecutor, guildID int64, iD int64, selectCols ...string) (*GenaiCommand, error) {
	genaiCommandObj := &GenaiCommand{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"genai_commands\" where \"guild_id\"=$1 AND \"id\"=$2", sel,
	)

	q := queries.Raw(query, guildID, iD)

	err := q.Bind(ctx, exec, genaiCommandObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from genai_commands")
	}

	return genaiCommandObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *GenaiCommand) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *GenaiCommand) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no genai_commands provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(genaiCommandColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	genaiCommandInsertCacheMut.RLock()
	cache, cached := genaiCommandInsertCache[key]
	genaiCommandInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			genaiCommandAllColumns,
			genaiCommandColumnsWithDefault,
			genaiCommandColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(genaiCommandType, genaiCommandMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(genaiCommandType, genaiCommandMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"genai_commands\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"genai_commands\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into genai_commands")
	}

	if !cached {
		genaiCommandInsertCacheMut.Lock()
		genaiCommandInsertCache[key] = cache
		genaiCommandInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single GenaiCommand record using the global executor.
// See Update for more documentation.
func (o *GenaiCommand) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the GenaiCommand.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *GenaiCommand) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	genaiCommandUpdateCacheMut.RLock()
	cache, cached := genaiCommandUpdateCache[key]
	genaiCommandUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			genaiCommandAllColumns,
			genaiCommandPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update genai_commands, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"genai_commands\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, genaiCommandPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(genaiCommandType, genaiCommandMapping, append(wl, genaiCommandPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update genai_commands row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for genai_commands")
	}

	if !cached {
		genaiCommandUpdateCacheMut.Lock()
		genaiCommandUpdateCache[key] = cache
		genaiCommandUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q genaiCommandQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q genaiCommandQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for genai_commands")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for genai_commands")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o GenaiCommandSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GenaiCommandSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), genaiCommandPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"genai_commands\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, genaiCommandPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in genaiCommand slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all genaiCommand")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *GenaiCommand) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *GenaiCommand) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no genai_commands provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(genaiCommandColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	genaiCommandUpsertCacheMut.RLock()
	cache, cached := genaiCommandUpsertCache[key]
	genaiCommandUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			genaiCommandAllColumns,
			genaiCommandColumnsWithDefault,
			genaiCommandColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			genaiCommandAllColumns,
			genaiCommandPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert genai_commands, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(genaiCommandPrimaryKeyColumns))
			copy(conflict, genaiCommandPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"genai_commands\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(genaiCommandType, genaiCommandMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(genaiCommandType, genaiCommandMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert genai_commands")
	}

	if !cached {
		genaiCommandUpsertCacheMut.Lock()
		genaiCommandUpsertCache[key] = cache
		genaiCommandUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single GenaiCommand record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *GenaiCommand) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single GenaiCommand record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *GenaiCommand) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no GenaiCommand provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), genaiCommandPrimaryKeyMapping)
	sql := "DELETE FROM \"genai_commands\" WHERE \"guild_id\"=$1 AND \"id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from genai_commands")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for genai_commands")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q genaiCommandQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no genaiCommandQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from genai_commands")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for genai_commands")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o GenaiCommandSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GenaiCommandSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), genaiCommandPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"genai_commands\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, genaiCommandPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from genaiCommand slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for genai_commands")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *GenaiCommand) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no GenaiCommand provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *GenaiCommand) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindGenaiCommand(ctx, exec, o.GuildID, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GenaiCommandSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty GenaiCommandSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GenaiCommandSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := GenaiCommandSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), genaiCommandPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"genai_commands\".* FROM \"genai_commands\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, genaiCommandPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GenaiCommandSlice")
	}

	*o = slice

	return nil
}

// GenaiCommandExistsG checks if the GenaiCommand row exists.
func GenaiCommandExistsG(ctx context.Context, guildID int64, iD int64) (bool, error) {
	return GenaiCommandExists(ctx, boil.GetContextDB(), guildID, iD)
}

// GenaiCommandExists checks if the GenaiCommand row exists.
func GenaiCommandExists(ctx context.Context, exec boil.ContextExecutor, guildID int64, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"genai_commands\" where \"guild_id\"=$1 AND \"id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, guildID, iD)
	}
	row := exec.QueryRowContext(ctx, sql, guildID, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if genai_commands exists")
	}

	return exists, nil
}
