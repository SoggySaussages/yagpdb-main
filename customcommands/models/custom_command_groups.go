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

// CustomCommandGroup is an object representing the database table.
type CustomCommandGroup struct {
	ID                int64            `boil:"id" json:"id" toml:"id" yaml:"id"`
	GuildID           int64            `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	Name              string           `boil:"name" json:"name" toml:"name" yaml:"name"`
	GitHub            string           `boil:"github" json:"github" toml:"github" yaml:"github"`
	IgnoreRoles       types.Int64Array `boil:"ignore_roles" json:"ignore_roles,omitempty" toml:"ignore_roles" yaml:"ignore_roles,omitempty"`
	IgnoreChannels    types.Int64Array `boil:"ignore_channels" json:"ignore_channels,omitempty" toml:"ignore_channels" yaml:"ignore_channels,omitempty"`
	WhitelistRoles    types.Int64Array `boil:"whitelist_roles" json:"whitelist_roles,omitempty" toml:"whitelist_roles" yaml:"whitelist_roles,omitempty"`
	WhitelistChannels types.Int64Array `boil:"whitelist_channels" json:"whitelist_channels,omitempty" toml:"whitelist_channels" yaml:"whitelist_channels,omitempty"`

	R *customCommandGroupR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L customCommandGroupL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CustomCommandGroupColumns = struct {
	ID                string
	GuildID           string
	Name              string
	GitHub            string
	IgnoreRoles       string
	IgnoreChannels    string
	WhitelistRoles    string
	WhitelistChannels string
}{
	ID:                "id",
	GuildID:           "guild_id",
	Name:              "name",
	GitHub:            "github",
	IgnoreRoles:       "ignore_roles",
	IgnoreChannels:    "ignore_channels",
	WhitelistRoles:    "whitelist_roles",
	WhitelistChannels: "whitelist_channels",
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

var CustomCommandGroupWhere = struct {
	ID                whereHelperint64
	GuildID           whereHelperint64
	Name              whereHelperstring
	GitHub            whereHelperstring
	IgnoreRoles       whereHelpertypes_Int64Array
	IgnoreChannels    whereHelpertypes_Int64Array
	WhitelistRoles    whereHelpertypes_Int64Array
	WhitelistChannels whereHelpertypes_Int64Array
}{
	ID:                whereHelperint64{field: "\"custom_command_groups\".\"id\""},
	GuildID:           whereHelperint64{field: "\"custom_command_groups\".\"guild_id\""},
	Name:              whereHelperstring{field: "\"custom_command_groups\".\"name\""},
	GitHub:            whereHelperstring{field: "\"custom_command_groups\".\"github\""},
	IgnoreRoles:       whereHelpertypes_Int64Array{field: "\"custom_command_groups\".\"ignore_roles\""},
	IgnoreChannels:    whereHelpertypes_Int64Array{field: "\"custom_command_groups\".\"ignore_channels\""},
	WhitelistRoles:    whereHelpertypes_Int64Array{field: "\"custom_command_groups\".\"whitelist_roles\""},
	WhitelistChannels: whereHelpertypes_Int64Array{field: "\"custom_command_groups\".\"whitelist_channels\""},
}

// CustomCommandGroupRels is where relationship names are stored.
var CustomCommandGroupRels = struct {
	GroupCustomCommands string
}{
	GroupCustomCommands: "GroupCustomCommands",
}

// customCommandGroupR is where relationships are stored.
type customCommandGroupR struct {
	GroupCustomCommands CustomCommandSlice
}

// NewStruct creates a new relationship struct
func (*customCommandGroupR) NewStruct() *customCommandGroupR {
	return &customCommandGroupR{}
}

// customCommandGroupL is where Load methods for each relationship are stored.
type customCommandGroupL struct{}

var (
	customCommandGroupAllColumns            = []string{"id", "guild_id", "name", "ignore_roles", "ignore_channels", "whitelist_roles", "whitelist_channels", "github"}
	customCommandGroupColumnsWithoutDefault = []string{"guild_id", "name", "ignore_roles", "ignore_channels", "whitelist_roles", "whitelist_channels"}
	customCommandGroupColumnsWithDefault    = []string{"id", "github"}
	customCommandGroupPrimaryKeyColumns     = []string{"id"}
)

type (
	// CustomCommandGroupSlice is an alias for a slice of pointers to CustomCommandGroup.
	// This should generally be used opposed to []CustomCommandGroup.
	CustomCommandGroupSlice []*CustomCommandGroup

	customCommandGroupQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	customCommandGroupType                 = reflect.TypeOf(&CustomCommandGroup{})
	customCommandGroupMapping              = queries.MakeStructMapping(customCommandGroupType)
	customCommandGroupPrimaryKeyMapping, _ = queries.BindMapping(customCommandGroupType, customCommandGroupMapping, customCommandGroupPrimaryKeyColumns)
	customCommandGroupInsertCacheMut       sync.RWMutex
	customCommandGroupInsertCache          = make(map[string]insertCache)
	customCommandGroupUpdateCacheMut       sync.RWMutex
	customCommandGroupUpdateCache          = make(map[string]updateCache)
	customCommandGroupUpsertCacheMut       sync.RWMutex
	customCommandGroupUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single customCommandGroup record from the query using the global executor.
func (q customCommandGroupQuery) OneG(ctx context.Context) (*CustomCommandGroup, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single customCommandGroup record from the query.
func (q customCommandGroupQuery) One(ctx context.Context, exec boil.ContextExecutor) (*CustomCommandGroup, error) {
	o := &CustomCommandGroup{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for custom_command_groups")
	}

	return o, nil
}

// AllG returns all CustomCommandGroup records from the query using the global executor.
func (q customCommandGroupQuery) AllG(ctx context.Context) (CustomCommandGroupSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all CustomCommandGroup records from the query.
func (q customCommandGroupQuery) All(ctx context.Context, exec boil.ContextExecutor) (CustomCommandGroupSlice, error) {
	var o []*CustomCommandGroup

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to CustomCommandGroup slice")
	}

	return o, nil
}

// CountG returns the count of all CustomCommandGroup records in the query, and panics on error.
func (q customCommandGroupQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all CustomCommandGroup records in the query.
func (q customCommandGroupQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count custom_command_groups rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q customCommandGroupQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q customCommandGroupQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if custom_command_groups exists")
	}

	return count > 0, nil
}

// GroupCustomCommands retrieves all the custom_command's CustomCommands with an executor via group_id column.
func (o *CustomCommandGroup) GroupCustomCommands(mods ...qm.QueryMod) customCommandQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"custom_commands\".\"group_id\"=?", o.ID),
	)

	query := CustomCommands(queryMods...)
	queries.SetFrom(query.Query, "\"custom_commands\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"custom_commands\".*"})
	}

	return query
}

// LoadGroupCustomCommands allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (customCommandGroupL) LoadGroupCustomCommands(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCustomCommandGroup interface{}, mods queries.Applicator) error {
	var slice []*CustomCommandGroup
	var object *CustomCommandGroup

	if singular {
		object = maybeCustomCommandGroup.(*CustomCommandGroup)
	} else {
		slice = *maybeCustomCommandGroup.(*[]*CustomCommandGroup)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &customCommandGroupR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &customCommandGroupR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`custom_commands`), qm.WhereIn(`custom_commands.group_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load custom_commands")
	}

	var resultSlice []*CustomCommand
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice custom_commands")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on custom_commands")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for custom_commands")
	}

	if singular {
		object.R.GroupCustomCommands = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &customCommandR{}
			}
			foreign.R.Group = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.GroupID) {
				local.R.GroupCustomCommands = append(local.R.GroupCustomCommands, foreign)
				if foreign.R == nil {
					foreign.R = &customCommandR{}
				}
				foreign.R.Group = local
				break
			}
		}
	}

	return nil
}

// AddGroupCustomCommandsG adds the given related objects to the existing relationships
// of the custom_command_group, optionally inserting them as new records.
// Appends related to o.R.GroupCustomCommands.
// Sets related.R.Group appropriately.
// Uses the global database handle.
func (o *CustomCommandGroup) AddGroupCustomCommandsG(ctx context.Context, insert bool, related ...*CustomCommand) error {
	return o.AddGroupCustomCommands(ctx, boil.GetContextDB(), insert, related...)
}

// AddGroupCustomCommands adds the given related objects to the existing relationships
// of the custom_command_group, optionally inserting them as new records.
// Appends related to o.R.GroupCustomCommands.
// Sets related.R.Group appropriately.
func (o *CustomCommandGroup) AddGroupCustomCommands(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*CustomCommand) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.GroupID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"custom_commands\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"group_id"}),
				strmangle.WhereClause("\"", "\"", 2, customCommandPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.GuildID, rel.LocalID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.GroupID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &customCommandGroupR{
			GroupCustomCommands: related,
		}
	} else {
		o.R.GroupCustomCommands = append(o.R.GroupCustomCommands, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &customCommandR{
				Group: o,
			}
		} else {
			rel.R.Group = o
		}
	}
	return nil
}

// SetGroupCustomCommandsG removes all previously related items of the
// custom_command_group replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Group's GroupCustomCommands accordingly.
// Replaces o.R.GroupCustomCommands with related.
// Sets related.R.Group's GroupCustomCommands accordingly.
// Uses the global database handle.
func (o *CustomCommandGroup) SetGroupCustomCommandsG(ctx context.Context, insert bool, related ...*CustomCommand) error {
	return o.SetGroupCustomCommands(ctx, boil.GetContextDB(), insert, related...)
}

// SetGroupCustomCommands removes all previously related items of the
// custom_command_group replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Group's GroupCustomCommands accordingly.
// Replaces o.R.GroupCustomCommands with related.
// Sets related.R.Group's GroupCustomCommands accordingly.
func (o *CustomCommandGroup) SetGroupCustomCommands(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*CustomCommand) error {
	query := "update \"custom_commands\" set \"group_id\" = null where \"group_id\" = $1"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.GroupCustomCommands {
			queries.SetScanner(&rel.GroupID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Group = nil
		}

		o.R.GroupCustomCommands = nil
	}
	return o.AddGroupCustomCommands(ctx, exec, insert, related...)
}

// RemoveGroupCustomCommandsG relationships from objects passed in.
// Removes related items from R.GroupCustomCommands (uses pointer comparison, removal does not keep order)
// Sets related.R.Group.
// Uses the global database handle.
func (o *CustomCommandGroup) RemoveGroupCustomCommandsG(ctx context.Context, related ...*CustomCommand) error {
	return o.RemoveGroupCustomCommands(ctx, boil.GetContextDB(), related...)
}

// RemoveGroupCustomCommands relationships from objects passed in.
// Removes related items from R.GroupCustomCommands (uses pointer comparison, removal does not keep order)
// Sets related.R.Group.
func (o *CustomCommandGroup) RemoveGroupCustomCommands(ctx context.Context, exec boil.ContextExecutor, related ...*CustomCommand) error {
	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.GroupID, nil)
		if rel.R != nil {
			rel.R.Group = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("group_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.GroupCustomCommands {
			if rel != ri {
				continue
			}

			ln := len(o.R.GroupCustomCommands)
			if ln > 1 && i < ln-1 {
				o.R.GroupCustomCommands[i] = o.R.GroupCustomCommands[ln-1]
			}
			o.R.GroupCustomCommands = o.R.GroupCustomCommands[:ln-1]
			break
		}
	}

	return nil
}

// CustomCommandGroups retrieves all the records using an executor.
func CustomCommandGroups(mods ...qm.QueryMod) customCommandGroupQuery {
	mods = append(mods, qm.From("\"custom_command_groups\""))
	return customCommandGroupQuery{NewQuery(mods...)}
}

// FindCustomCommandGroupG retrieves a single record by ID.
func FindCustomCommandGroupG(ctx context.Context, iD int64, selectCols ...string) (*CustomCommandGroup, error) {
	return FindCustomCommandGroup(ctx, boil.GetContextDB(), iD, selectCols...)
}

// FindCustomCommandGroup retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCustomCommandGroup(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*CustomCommandGroup, error) {
	customCommandGroupObj := &CustomCommandGroup{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"custom_command_groups\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, customCommandGroupObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from custom_command_groups")
	}

	return customCommandGroupObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *CustomCommandGroup) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *CustomCommandGroup) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no custom_command_groups provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(customCommandGroupColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	customCommandGroupInsertCacheMut.RLock()
	cache, cached := customCommandGroupInsertCache[key]
	customCommandGroupInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			customCommandGroupAllColumns,
			customCommandGroupColumnsWithDefault,
			customCommandGroupColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(customCommandGroupType, customCommandGroupMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(customCommandGroupType, customCommandGroupMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"custom_command_groups\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"custom_command_groups\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into custom_command_groups")
	}

	if !cached {
		customCommandGroupInsertCacheMut.Lock()
		customCommandGroupInsertCache[key] = cache
		customCommandGroupInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single CustomCommandGroup record using the global executor.
// See Update for more documentation.
func (o *CustomCommandGroup) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the CustomCommandGroup.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *CustomCommandGroup) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	customCommandGroupUpdateCacheMut.RLock()
	cache, cached := customCommandGroupUpdateCache[key]
	customCommandGroupUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			customCommandGroupAllColumns,
			customCommandGroupPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update custom_command_groups, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"custom_command_groups\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, customCommandGroupPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(customCommandGroupType, customCommandGroupMapping, append(wl, customCommandGroupPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update custom_command_groups row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for custom_command_groups")
	}

	if !cached {
		customCommandGroupUpdateCacheMut.Lock()
		customCommandGroupUpdateCache[key] = cache
		customCommandGroupUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q customCommandGroupQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q customCommandGroupQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for custom_command_groups")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for custom_command_groups")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CustomCommandGroupSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CustomCommandGroupSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customCommandGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"custom_command_groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, customCommandGroupPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in customCommandGroup slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all customCommandGroup")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *CustomCommandGroup) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *CustomCommandGroup) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no custom_command_groups provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(customCommandGroupColumnsWithDefault, o)

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

	customCommandGroupUpsertCacheMut.RLock()
	cache, cached := customCommandGroupUpsertCache[key]
	customCommandGroupUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			customCommandGroupAllColumns,
			customCommandGroupColumnsWithDefault,
			customCommandGroupColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			customCommandGroupAllColumns,
			customCommandGroupPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert custom_command_groups, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(customCommandGroupPrimaryKeyColumns))
			copy(conflict, customCommandGroupPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"custom_command_groups\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(customCommandGroupType, customCommandGroupMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(customCommandGroupType, customCommandGroupMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert custom_command_groups")
	}

	if !cached {
		customCommandGroupUpsertCacheMut.Lock()
		customCommandGroupUpsertCache[key] = cache
		customCommandGroupUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single CustomCommandGroup record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *CustomCommandGroup) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single CustomCommandGroup record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *CustomCommandGroup) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no CustomCommandGroup provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), customCommandGroupPrimaryKeyMapping)
	sql := "DELETE FROM \"custom_command_groups\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from custom_command_groups")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for custom_command_groups")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q customCommandGroupQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no customCommandGroupQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from custom_command_groups")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for custom_command_groups")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o CustomCommandGroupSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CustomCommandGroupSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customCommandGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"custom_command_groups\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, customCommandGroupPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from customCommandGroup slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for custom_command_groups")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *CustomCommandGroup) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no CustomCommandGroup provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *CustomCommandGroup) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCustomCommandGroup(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CustomCommandGroupSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty CustomCommandGroupSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CustomCommandGroupSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CustomCommandGroupSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customCommandGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"custom_command_groups\".* FROM \"custom_command_groups\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, customCommandGroupPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in CustomCommandGroupSlice")
	}

	*o = slice

	return nil
}

// CustomCommandGroupExistsG checks if the CustomCommandGroup row exists.
func CustomCommandGroupExistsG(ctx context.Context, iD int64) (bool, error) {
	return CustomCommandGroupExists(ctx, boil.GetContextDB(), iD)
}

// CustomCommandGroupExists checks if the CustomCommandGroup row exists.
func CustomCommandGroupExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"custom_command_groups\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if custom_command_groups exists")
	}

	return exists, nil
}
