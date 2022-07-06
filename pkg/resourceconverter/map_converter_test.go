package resourceconverter

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestMapConverter(t *testing.T) {
	ctx := context.Background()
	t.Run("TagsPassedIn", func(t *testing.T) {
		entry := map[string]any{
			"ID":        "id1",
			"Attr1":     1,
			"Attr2":     "hi",
			"Attr3":     map[string]interface{}{"a": "b", "c": 2},
			"WeirdTags": []WeirdTags{{WeirdKey: "key1", WeirdValue: "val1"}, {WeirdKey: "key2", WeirdValue: "val2"}},
		}
		rC := &MapConverter{
			IdField:      "ID",
			ResourceType: "DummyResource",
		}
		resource := model.Resource{
			Region: "dummyRegion",
			Type:   "DummyResource",
			Tags:   model.Tags{{Key: "key1", Value: "val3"}, {Key: "key2", Value: "val4"}},
		}
		err := rC.ToResource(ctx, entry, &resource)
		require.NoError(t, err)
		expectedResource := model.Resource{
			Region:  "dummyRegion",
			Id:      "id1",
			Type:    "DummyResource",
			Tags:    model.Tags{{Key: "key1", Value: "val3"}, {Key: "key2", Value: "val4"}},
			RawData: datatypes.JSON([]byte(`{"ID":"id1","Attr1":1,"Attr2":"hi","Attr3":{"a":"b","c":2},"WeirdTags":[{"WeirdKey":"key1","WeirdValue":"val1"},{"WeirdKey":"key2","WeirdValue":"val2"}]}`)),
		}
		testingutil.AssertEqualsResource(t, expectedResource, resource)
	})
}
