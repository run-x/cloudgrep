package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_ec2(mapping map[string]mapper) {
	mapping["ec2.Instance"] = mapper{
		FetchFunc: p.fetch_ec2_Instance,
		IdField:   "InstanceId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.VPC"] = mapper{
		FetchFunc: p.fetch_ec2_VPC,
		IdField:   "VpcId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Volume"] = mapper{
		FetchFunc: p.fetch_ec2_Volume,
		IdField:   "VolumeId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_ec2_Instance(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeInstancesInput{}
	input.Filters = describeInstancesFilters()

	resourceConverter := p.converterFor("ec2.Instance")
	paginator := ec2.NewDescribeInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Instance", err)
		}

		for _, item_0 := range page.Reservations {
			if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, item_0.Instances); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_VPC(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVpcsInput{}

	resourceConverter := p.converterFor("ec2.VPC")
	paginator := ec2.NewDescribeVpcsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.VPC", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Vpcs); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Volume(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVolumesInput{}

	resourceConverter := p.converterFor("ec2.Volume")
	paginator := ec2.NewDescribeVolumesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Volume", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.Volumes); err != nil {
			return err
		}
	}

	return nil
}