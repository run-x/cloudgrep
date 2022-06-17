package aws

import (
	"fmt"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchCloudfrontDistributions(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "cloudfront.Distribution")

	testingutil.AssertResourceCount(t, resources, "", 1)
	fmt.Println(resources[0])
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type: "cloudfront.Distribution" +
			"",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "cloudfront-distribution-0",
			},
		},
	})
}
