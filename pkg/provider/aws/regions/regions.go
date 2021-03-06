package regions

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"golang.org/x/exp/slices"
)

const Global = "global"
const All = "all"

// SelectRegions returns the regions the user has selected, from either the cloudgrep config, AWS config, or prompting.
// The special value "all" can be present by itself in configuredRegions to automatically select all enabled regions in the account.
func SelectRegions(ctx context.Context, configuredRegions []string, awsConfig aws.Config) ([]Region, error) {
	var err error

	if len(configuredRegions) == 1 && configuredRegions[0] == All {
		return allRegions(ctx, awsConfig)
	}

	if slices.Contains(configuredRegions, All) {
		return nil, fmt.Errorf("can only use '%s' as a region if it is the only configured region", All)
	}

	if len(configuredRegions) > 0 {
		// If regions were configured, use those
		err = validateRegions(configuredRegions)
		if err != nil {
			return nil, fmt.Errorf("unable to use configured regions: %w", err)
		}

		return regionsFromStrings(configuredRegions)
	}

	region := awsConfig.Region

	// If we can't detect region automatically, prompt for it
	if region == "" {
		region, err = promptForRegion(ctx)
		if err != nil {
			if err == ctx.Err() {
				return nil, err
			}

			return nil, fmt.Errorf("error prompting for region: %w", err)
		}
	} else {
		err = validateRegions([]string{region})
		if err != nil {
			return nil, err
		}
	}

	if region == All {
		return allRegions(ctx, awsConfig)
	}

	// Always include global region without explicit configuration excluding it
	regions := []string{Global, region}

	return regionsFromStrings(regions)
}

// IsValid returns true if the given region is recognized as valid.
func IsValid(region string) bool {
	if region == Global || region == All {
		return true
	}

	_, has := officialRegions[region]
	return has
}

// SetConfigRegion updates the aws.Config value with one of the regions in the passed list, to ensure
// there is always configured region.
func SetConfigRegion(cfg *aws.Config, regions []Region) {
	if cfg == nil {
		panic("unexpected nil cfg")
	}

	if cfg.Region != "" {
		return
	}

	cfg.Region = "us-east-1"
	for _, region := range regions {
		if !region.IsGlobal() {
			cfg.Region = region.ID()
			return
		}
	}
}
