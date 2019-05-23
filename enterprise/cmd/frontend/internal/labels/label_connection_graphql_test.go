package labels

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/db"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/types"
)

func TestGraphQL_LabelConnection(t *testing.T) {
	resetMocks()
	const (
		wantThreadID = 3
		wantLabelID  = 2
	)
	db.Mocks.DiscussionThreads.Get = func(int64) (*types.DiscussionThread, error) {
		return &types.DiscussionThread{ID: wantThreadID}, nil
	}
	mocks.labelsObjects.List = func(dbLabelsObjectsListOptions) ([]*dbObjectLabel, error) {
		return []*dbObjectLabel{{Thread: wantThreadID, Label: wantLabelID}}, nil
	}
	mocks.labels.GetByID = func(id int64) (*dbLabel, error) {
		if id != wantLabelID {
			t.Errorf("got %d, want %d", id, wantLabelID)
		}
		return &dbLabel{Name: "n"}, nil
	}

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Context: backend.WithAuthzBypass(context.Background()),
			Schema:  graphqlbackend.GraphQLSchema,
			Query: `
				{
					node(id: "RGlzY3Vzc2lvblRocmVhZDoiMyI=") {
						... on DiscussionThread {
							labels {
								nodes {
									name
								}
								totalCount
								pageInfo {
									hasNextPage
								}
							}
						}
					}
				}
			`,
			ExpectedResult: `
				{
					"node": {
						"labels": {
							"nodes": [
								{
									"name": "n"
								}
							],
							"totalCount": 1,
							"pageInfo": {
								"hasNextPage": false
							}
						}
					}
				}
			`,
		},
	})
}
