package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_rds(mapping map[string]mapper) {
	mapping["rds.DBCluster"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetch_rds_DBCluster,
		IdField:           "DBClusterIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["rds.DBInstance"] = mapper{
		ServiceEndpointID: "rds",
		FetchFunc:         p.fetch_rds_DBInstance,
		IdField:           "DBInstanceIdentifier",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagList",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_rds_DBCluster(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBClustersInput{}

	resourceConverter := p.converterFor("rds.DBCluster")
	paginator := rds.NewDescribeDBClustersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBCluster", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.DBClusters); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_rds_DBInstance(ctx context.Context, output chan<- model.Resource) error {
	client := rds.NewFromConfig(p.config)
	input := &rds.DescribeDBInstancesInput{}

	resourceConverter := p.converterFor("rds.DBInstance")
	paginator := rds.NewDescribeDBInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "rds.DBInstance", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.DBInstances); err != nil {
			return err
		}
	}

	return nil
}
