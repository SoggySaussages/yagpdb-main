// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testRedditFeeds(t *testing.T) {
	t.Parallel()

	query := RedditFeeds()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testRedditFeedsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRedditFeedsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := RedditFeeds().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRedditFeedsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RedditFeedSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRedditFeedsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := RedditFeedExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if RedditFeed exists: %s", err)
	}
	if !e {
		t.Errorf("Expected RedditFeedExists to return true, but got false.")
	}
}

func testRedditFeedsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	redditFeedFound, err := FindRedditFeed(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if redditFeedFound == nil {
		t.Error("want a record, got nil")
	}
}

func testRedditFeedsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = RedditFeeds().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testRedditFeedsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := RedditFeeds().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testRedditFeedsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	redditFeedOne := &RedditFeed{}
	redditFeedTwo := &RedditFeed{}
	if err = randomize.Struct(seed, redditFeedOne, redditFeedDBTypes, false, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}
	if err = randomize.Struct(seed, redditFeedTwo, redditFeedDBTypes, false, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = redditFeedOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = redditFeedTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RedditFeeds().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testRedditFeedsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	redditFeedOne := &RedditFeed{}
	redditFeedTwo := &RedditFeed{}
	if err = randomize.Struct(seed, redditFeedOne, redditFeedDBTypes, false, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}
	if err = randomize.Struct(seed, redditFeedTwo, redditFeedDBTypes, false, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = redditFeedOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = redditFeedTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testRedditFeedsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRedditFeedsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(redditFeedColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRedditFeedsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testRedditFeedsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RedditFeedSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testRedditFeedsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RedditFeeds().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	redditFeedDBTypes = map[string]string{`ID`: `bigint`, `GuildID`: `bigint`, `ChannelID`: `bigint`, `Subreddit`: `text`, `FilterNSFW`: `integer`, `MinUpvotes`: `integer`, `UseEmbeds`: `boolean`, `Slow`: `boolean`, `Disabled`: `boolean`}
	_                 = bytes.MinRead
)

func testRedditFeedsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(redditFeedPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(redditFeedColumns) == len(redditFeedPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testRedditFeedsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(redditFeedColumns) == len(redditFeedPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RedditFeed{}
	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, redditFeedDBTypes, true, redditFeedPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(redditFeedColumns, redditFeedPrimaryKeyColumns) {
		fields = redditFeedColumns
	} else {
		fields = strmangle.SetComplement(
			redditFeedColumns,
			redditFeedPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := RedditFeedSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testRedditFeedsUpsert(t *testing.T) {
	t.Parallel()

	if len(redditFeedColumns) == len(redditFeedPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := RedditFeed{}
	if err = randomize.Struct(seed, &o, redditFeedDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RedditFeed: %s", err)
	}

	count, err := RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, redditFeedDBTypes, false, redditFeedPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RedditFeed struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RedditFeed: %s", err)
	}

	count, err = RedditFeeds().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}