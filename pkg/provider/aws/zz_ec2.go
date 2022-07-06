package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_ec2(mapping map[string]mapper) {
	mapping["ec2.Address"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Address,
		IdField:           "AllocationId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.CapacityReservation"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_CapacityReservation,
		IdField:           "CapacityReservationId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ClientVpnEndpoint"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_ClientVpnEndpoint,
		IdField:           "ClientVpnEndpointId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Fleet"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Fleet,
		IdField:           "FleetId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.FlowLog"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_FlowLog,
		IdField:           "FlowLogId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Image"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Image,
		IdField:           "ImageId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Instance"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Instance,
		IdField:           "InstanceId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.KeyPair"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_KeyPair,
		IdField:           "KeyPairId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.LaunchTemplate"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_LaunchTemplate,
		IdField:           "LaunchTemplateId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NatGateway"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_NatGateway,
		IdField:           "NatGatewayId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NetworkAcl"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_NetworkAcl,
		IdField:           "NetworkAclId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.NetworkInterface"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_NetworkInterface,
		IdField:           "NetworkInterfaceId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "TagSet",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.ReservedInstance"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_ReservedInstance,
		IdField:           "ReservedInstancesId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.RouteTable"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_RouteTable,
		IdField:           "RouteTableId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.SecurityGroup"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_SecurityGroup,
		IdField:           "GroupId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Snapshot"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Snapshot,
		IdField:           "SnapshotId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.SpotInstanceRequest"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_SpotInstanceRequest,
		IdField:           "SpotInstanceRequestId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Subnet"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Subnet,
		IdField:           "SubnetId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Volume"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Volume,
		IdField:           "VolumeId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
	mapping["ec2.Vpc"] = mapper{
		ServiceEndpointID: "ec2",
		FetchFunc:         p.fetch_ec2_Vpc,
		IdField:           "VpcId",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_ec2_Address(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeAddressesInput{}

	commonTransformers := p.baseTransformers("ec2.Address")
	converter := p.converterFor("ec2.Address")
	transformers := append(
		resourceconverter.AllToGeneric[types.Address](commonTransformers...),
		resourceconverter.WithConverter[types.Address](converter),
	)
	results, err := client.DescribeAddresses(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.Address", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.Addresses, transformers...); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetch_ec2_CapacityReservation(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeCapacityReservationsInput{}
	input.Filters = describeCapacityReservationsFilters()

	commonTransformers := p.baseTransformers("ec2.CapacityReservation")
	converter := p.converterFor("ec2.CapacityReservation")
	transformers := append(
		resourceconverter.AllToGeneric[types.CapacityReservation](commonTransformers...),
		resourceconverter.WithConverter[types.CapacityReservation](converter),
	)
	paginator := ec2.NewDescribeCapacityReservationsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.CapacityReservation", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.CapacityReservations, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_ClientVpnEndpoint(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeClientVpnEndpointsInput{}

	commonTransformers := p.baseTransformers("ec2.ClientVpnEndpoint")
	converter := p.converterFor("ec2.ClientVpnEndpoint")
	transformers := append(
		resourceconverter.AllToGeneric[types.ClientVpnEndpoint](commonTransformers...),
		resourceconverter.WithConverter[types.ClientVpnEndpoint](converter),
	)
	paginator := ec2.NewDescribeClientVpnEndpointsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.ClientVpnEndpoint", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.ClientVpnEndpoints, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Fleet(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeFleetsInput{}
	input.Filters = describeFleetsFilters()

	commonTransformers := p.baseTransformers("ec2.Fleet")
	converter := p.converterFor("ec2.Fleet")
	transformers := append(
		resourceconverter.AllToGeneric[types.FleetData](commonTransformers...),
		resourceconverter.WithConverter[types.FleetData](converter),
	)
	paginator := ec2.NewDescribeFleetsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Fleet", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.Fleets, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_FlowLog(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeFlowLogsInput{}

	commonTransformers := p.baseTransformers("ec2.FlowLog")
	converter := p.converterFor("ec2.FlowLog")
	transformers := append(
		resourceconverter.AllToGeneric[types.FlowLog](commonTransformers...),
		resourceconverter.WithConverter[types.FlowLog](converter),
	)
	paginator := ec2.NewDescribeFlowLogsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.FlowLog", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.FlowLogs, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Image(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeImagesInput{}
	input.Owners = describeImagesOwners()

	commonTransformers := p.baseTransformers("ec2.Image")
	converter := p.converterFor("ec2.Image")
	transformers := append(
		resourceconverter.AllToGeneric[types.Image](commonTransformers...),
		resourceconverter.WithConverter[types.Image](converter),
	)
	results, err := client.DescribeImages(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.Image", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.Images, transformers...); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetch_ec2_Instance(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeInstancesInput{}
	input.Filters = describeInstancesFilters()

	commonTransformers := p.baseTransformers("ec2.Instance")
	converter := p.converterFor("ec2.Instance")
	transformers := append(
		resourceconverter.AllToGeneric[types.Instance](commonTransformers...),
		resourceconverter.WithConverter[types.Instance](converter),
	)
	paginator := ec2.NewDescribeInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Instance", err)
		}

		for _, item_0 := range page.Reservations {
			if err := resourceconverter.SendAll(ctx, output, item_0.Instances, transformers...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_KeyPair(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeKeyPairsInput{}

	commonTransformers := p.baseTransformers("ec2.KeyPair")
	converter := p.converterFor("ec2.KeyPair")
	transformers := append(
		resourceconverter.AllToGeneric[types.KeyPairInfo](commonTransformers...),
		resourceconverter.WithConverter[types.KeyPairInfo](converter),
	)
	results, err := client.DescribeKeyPairs(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.KeyPair", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.KeyPairs, transformers...); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetch_ec2_LaunchTemplate(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeLaunchTemplatesInput{}

	commonTransformers := p.baseTransformers("ec2.LaunchTemplate")
	converter := p.converterFor("ec2.LaunchTemplate")
	transformers := append(
		resourceconverter.AllToGeneric[types.LaunchTemplate](commonTransformers...),
		resourceconverter.WithConverter[types.LaunchTemplate](converter),
	)
	paginator := ec2.NewDescribeLaunchTemplatesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.LaunchTemplate", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.LaunchTemplates, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_NatGateway(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeNatGatewaysInput{}

	commonTransformers := p.baseTransformers("ec2.NatGateway")
	converter := p.converterFor("ec2.NatGateway")
	transformers := append(
		resourceconverter.AllToGeneric[types.NatGateway](commonTransformers...),
		resourceconverter.WithConverter[types.NatGateway](converter),
	)
	paginator := ec2.NewDescribeNatGatewaysPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.NatGateway", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.NatGateways, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_NetworkAcl(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeNetworkAclsInput{}

	commonTransformers := p.baseTransformers("ec2.NetworkAcl")
	converter := p.converterFor("ec2.NetworkAcl")
	transformers := append(
		resourceconverter.AllToGeneric[types.NetworkAcl](commonTransformers...),
		resourceconverter.WithConverter[types.NetworkAcl](converter),
	)
	paginator := ec2.NewDescribeNetworkAclsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.NetworkAcl", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.NetworkAcls, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_NetworkInterface(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeNetworkInterfacesInput{}

	commonTransformers := p.baseTransformers("ec2.NetworkInterface")
	converter := p.converterFor("ec2.NetworkInterface")
	transformers := append(
		resourceconverter.AllToGeneric[types.NetworkInterface](commonTransformers...),
		resourceconverter.WithConverter[types.NetworkInterface](converter),
	)
	paginator := ec2.NewDescribeNetworkInterfacesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.NetworkInterface", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.NetworkInterfaces, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_ReservedInstance(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeReservedInstancesInput{}
	input.Filters = describeReservedInstancesFilters()

	commonTransformers := p.baseTransformers("ec2.ReservedInstance")
	converter := p.converterFor("ec2.ReservedInstance")
	transformers := append(
		resourceconverter.AllToGeneric[types.ReservedInstances](commonTransformers...),
		resourceconverter.WithConverter[types.ReservedInstances](converter),
	)
	results, err := client.DescribeReservedInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "ec2.ReservedInstance", err)
	}
	if err := resourceconverter.SendAll(ctx, output, results.ReservedInstances, transformers...); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fetch_ec2_RouteTable(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeRouteTablesInput{}

	commonTransformers := p.baseTransformers("ec2.RouteTable")
	converter := p.converterFor("ec2.RouteTable")
	transformers := append(
		resourceconverter.AllToGeneric[types.RouteTable](commonTransformers...),
		resourceconverter.WithConverter[types.RouteTable](converter),
	)
	paginator := ec2.NewDescribeRouteTablesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.RouteTable", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.RouteTables, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_SecurityGroup(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSecurityGroupsInput{}

	commonTransformers := p.baseTransformers("ec2.SecurityGroup")
	converter := p.converterFor("ec2.SecurityGroup")
	transformers := append(
		resourceconverter.AllToGeneric[types.SecurityGroup](commonTransformers...),
		resourceconverter.WithConverter[types.SecurityGroup](converter),
	)
	paginator := ec2.NewDescribeSecurityGroupsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.SecurityGroup", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.SecurityGroups, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Snapshot(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSnapshotsInput{}
	input.OwnerIds = describeSnapshotsOwners()

	commonTransformers := p.baseTransformers("ec2.Snapshot")
	converter := p.converterFor("ec2.Snapshot")
	transformers := append(
		resourceconverter.AllToGeneric[types.Snapshot](commonTransformers...),
		resourceconverter.WithConverter[types.Snapshot](converter),
	)
	paginator := ec2.NewDescribeSnapshotsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Snapshot", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.Snapshots, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_SpotInstanceRequest(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSpotInstanceRequestsInput{}
	input.Filters = describeSpotInstanceRequestsFilters()

	commonTransformers := p.baseTransformers("ec2.SpotInstanceRequest")
	converter := p.converterFor("ec2.SpotInstanceRequest")
	transformers := append(
		resourceconverter.AllToGeneric[types.SpotInstanceRequest](commonTransformers...),
		resourceconverter.WithConverter[types.SpotInstanceRequest](converter),
	)
	paginator := ec2.NewDescribeSpotInstanceRequestsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.SpotInstanceRequest", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.SpotInstanceRequests, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Subnet(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeSubnetsInput{}

	commonTransformers := p.baseTransformers("ec2.Subnet")
	converter := p.converterFor("ec2.Subnet")
	transformers := append(
		resourceconverter.AllToGeneric[types.Subnet](commonTransformers...),
		resourceconverter.WithConverter[types.Subnet](converter),
	)
	paginator := ec2.NewDescribeSubnetsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Subnet", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.Subnets, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Volume(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVolumesInput{}

	commonTransformers := p.baseTransformers("ec2.Volume")
	converter := p.converterFor("ec2.Volume")
	transformers := append(
		resourceconverter.AllToGeneric[types.Volume](commonTransformers...),
		resourceconverter.WithConverter[types.Volume](converter),
	)
	paginator := ec2.NewDescribeVolumesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Volume", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.Volumes, transformers...); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) fetch_ec2_Vpc(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)
	input := &ec2.DescribeVpcsInput{}

	commonTransformers := p.baseTransformers("ec2.Vpc")
	converter := p.converterFor("ec2.Vpc")
	transformers := append(
		resourceconverter.AllToGeneric[types.Vpc](commonTransformers...),
		resourceconverter.WithConverter[types.Vpc](converter),
	)
	paginator := ec2.NewDescribeVpcsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "ec2.Vpc", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.Vpcs, transformers...); err != nil {
			return err
		}
	}

	return nil
}
