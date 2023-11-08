// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build e2e

package aggregator

import (
	"context"
	"errors"
	"testing"

	"github.com/chenmingyong0423/go-mongox/converter"

	"github.com/chenmingyong0423/go-mongox/builder/aggregation"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func getCollection(t *testing.T) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username:   "test",
		Password:   "test",
		AuthSource: "db-test",
	}))
	assert.NoError(t, err)
	assert.NoError(t, client.Ping(context.Background(), readpref.Primary()))

	return client.Database("db-test").Collection("test_user")
}

func TestAggregator_e2e_New(t *testing.T) {
	collection := getCollection(t)

	result := NewAggregator[types.TestUser](collection)
	assert.NotNil(t, result, "Expected non-nil Aggregator")
	assert.Equal(t, collection, result.collection, "Expected collection field to be initialized correctly")
}

func TestAggregator_e2e_Aggregation(t *testing.T) {
	collection := getCollection(t)
	aggregator := NewAggregator[types.TestUser](collection)

	testCases := []struct {
		name   string
		before func(ctx context.Context, t *testing.T)
		after  func(ctx context.Context, t *testing.T)

		pipeline           any
		aggregationOptions []*options.AggregateOptions

		ctx     context.Context
		want    []*types.TestUser
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "got error when pipeline is nil",
			before: func(_ context.Context, _ *testing.T) {},
			after:  func(_ context.Context, _ *testing.T) {},

			pipeline:           nil,
			aggregationOptions: nil,

			ctx:     context.Background(),
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "decode error",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					&types.IllegalUser{
						Id: "1", Name: "cmy", Age: "24",
					},
					&types.IllegalUser{
						Id: "2", Name: "gopher", Age: "20",
					},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []string{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", "1", "2").Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline:           mongo.Pipeline{},
			aggregationOptions: nil,
			want:               []*types.TestUser{},
			ctx:                context.Background(),
			wantErr:            assert.Error,
		},
		{
			name: "got result when pipeline is empty",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					types.TestUser{Id: "1", Name: "cmy", Age: 24},
					types.TestUser{Id: "2", Name: "gopher", Age: 20},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []any{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", []string{"1", "2"}...).Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline:           mongo.Pipeline{},
			aggregationOptions: nil,
			want: []*types.TestUser{
				{Id: "2", Name: "gopher", Age: 20},
				{Id: "1", Name: "cmy", Age: 24},
			},
			ctx:     context.Background(),
			wantErr: assert.NoError,
		},
		{
			name: "got result by pipeline with match stage",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					types.TestUser{Id: "2", Name: "gopher", Age: 20},
					types.TestUser{Id: "1", Name: "cmy", Age: 24},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []any{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", []string{"1", "2"}...).Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline: aggregation.StageBsonBuilder().Sort(converter.KeyValue[any]("age", -1)).Build(),
			want: []*types.TestUser{
				{Id: "1", Name: "cmy", Age: 24},
				{Id: "2", Name: "gopher", Age: 20},
			},
			ctx:     context.Background(),
			wantErr: assert.NoError,
		},
		{
			name: "got result with aggregation options",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					types.TestUser{Id: "2", Name: "gopher", Age: 20},
					types.TestUser{Id: "1", Name: "cmy", Age: 24},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []any{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", []string{"1", "2"}...).Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline: aggregation.StageBsonBuilder().Sort(converter.KeyValue[any]("name", 1)).Build(),
			aggregationOptions: []*options.AggregateOptions{
				options.Aggregate().SetCollation(&options.Collation{Locale: "en", Strength: 2}),
			},
			want: []*types.TestUser{
				{Id: "1", Name: "cmy", Age: 24},
				{Id: "2", Name: "gopher", Age: 20},
			},
			ctx:     context.Background(),
			wantErr: assert.NoError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(tc.ctx, t)
			testUsers, err := aggregator.Pipeline(tc.pipeline).AggregateOptions(tc.aggregationOptions...).Aggregation(tc.ctx)
			tc.after(tc.ctx, t)
			if tc.wantErr(t, err) {
				assert.ElementsMatch(t, tc.want, testUsers)
			}
		})
	}
}

func TestAggregator_e2e_AggregationWithCallback(t *testing.T) {
	collection := getCollection(t)
	aggregator := NewAggregator[types.TestUser](collection)

	type User struct {
		Id           string `bson:"_id"`
		Name         string `bson:"name"`
		Age          int64
		IsProgrammer bool `bson:"is_programmer"`
	}

	testCases := []struct {
		name               string
		before             func(ctx context.Context, t *testing.T)
		after              func(ctx context.Context, t *testing.T)
		pipeline           any
		aggregationOptions *options.AggregateOptions
		ctx                context.Context
		preUsers           []*User
		callback           types.ResultHandler
		want               []*User
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name:   "got error when pipeline is nil",
			before: func(_ context.Context, _ *testing.T) {},
			after:  func(_ context.Context, _ *testing.T) {},

			pipeline:           nil,
			aggregationOptions: nil,
			ctx:                context.Background(),
			want:               nil,
			wantErr:            assert.Error,
		},
		{
			name: "got result by pipeline with match stage",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					types.TestUser{Id: "2", Name: "gopher", Age: 20},
					types.TestUser{Id: "1", Name: "cmy", Age: 24},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []any{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", []string{"1", "2"}...).Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline: aggregation.StageBsonBuilder().Set(converter.KeyValue[any]("is_programmer", true)).Build(),
			preUsers: make([]*User, 0, 4),
			want: []*User{
				{Id: "1", Name: "cmy", Age: 24, IsProgrammer: true},
				{Id: "2", Name: "gopher", Age: 20, IsProgrammer: true},
			},
			ctx:     context.Background(),
			wantErr: assert.NoError,
		},
		{
			name: "got result with aggregation options",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					types.TestUser{Id: "2", Name: "gopher", Age: 20},
					types.TestUser{Id: "1", Name: "cmy", Age: 24},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []any{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", []string{"1", "2"}...).Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline: aggregation.StageBsonBuilder().Set(converter.KeyValue[any]("is_programmer", true)).Sort(converter.KeyValue[any]("name", 1)).Build(),
			preUsers: make([]*User, 0, 4),
			want: []*User{
				{Id: "1", Name: "cmy", Age: 24, IsProgrammer: true},
				{Id: "2", Name: "gopher", Age: 20, IsProgrammer: true},
			},
			aggregationOptions: options.Aggregate().SetCollation(&options.Collation{Locale: "en", Strength: 2}),
			ctx:                context.Background(),
			wantErr:            assert.NoError,
		},
		{
			name: "got error from cursor",
			before: func(ctx context.Context, t *testing.T) {
				insertManyResult, err := collection.InsertMany(ctx, []any{
					types.TestUser{Id: "2", Name: "gopher", Age: 20},
					types.TestUser{Id: "1", Name: "cmy", Age: 24},
				})
				assert.NoError(t, err)
				assert.ElementsMatch(t, []any{"1", "2"}, insertManyResult.InsertedIDs)
			},
			after: func(ctx context.Context, t *testing.T) {
				deleteResult, err := collection.DeleteMany(ctx, query.BsonBuilder().InString("_id", []string{"1", "2"}...).Build())
				assert.NoError(t, err)
				assert.Equal(t, int64(2), deleteResult.DeletedCount)
			},
			pipeline: aggregation.StageBsonBuilder().Set(converter.KeyValue[any]("is_programmer", true)).Sort(converter.KeyValue[any]("name", 1)).Build(),
			preUsers: make([]*User, 0),
			callback: func(cursor *mongo.Cursor) error {
				return errors.New("got error from cursor")
			},
			want:               []*User{},
			aggregationOptions: options.Aggregate().SetCollation(&options.Collation{Locale: "en", Strength: 2}),
			ctx:                context.Background(),
			wantErr:            assert.Error,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(tc.ctx, t)
			callback := func(cursor *mongo.Cursor) error {
				return cursor.All(context.Background(), &tc.preUsers)
			}
			if tc.callback != nil {
				callback = tc.callback
			}
			err := aggregator.Pipeline(tc.pipeline).AggregateOptions(tc.aggregationOptions).AggregationWithCallback(tc.ctx, callback)
			tc.after(tc.ctx, t)
			if tc.wantErr(t, err) {
				assert.ElementsMatch(t, tc.want, tc.preUsers)
			}
		})
	}
}
