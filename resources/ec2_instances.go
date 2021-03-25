package resources

import (
	"context"
	"github.com/cloudquery/cloudquery-plugin-sdk/plugin/schema"
)

func Ec2instances() *schema.Table {
	return &schema.Table{
		Name:     "aws_ec2_instances",
		Resolver: FetchEc2Instances,
		Columns: []schema.Column{
			{
				Name: "ami_launch_index",
				Type: schema.TypeInt,
			},
			{
				Name: "architecture",
				Type: schema.TypeString,
			},
			{
				Name: "capacity_reservation_id",
				Type: schema.TypeString,
			},
			{
				Name: "capacity_reservation_specification",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "capacity_reservation_specification_capacity_reservation_preference",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("CapacityReservationSpecification.CapacityReservationPreference"),
					},
					{
						Name:     "capacity_reservation_specification_capacity_reservation_target",
						Type:     schema.TypeEmbedded,
						Resolver: schema.PathResolver("CapacityReservationSpecification.CapacityReservationTarget"),
					},
				},
			},
			{
				Name: "client_token",
				Type: schema.TypeString,
			},
			{
				Name: "cpu_options",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "cpu_options_core_count",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("CpuOptions.CoreCount"),
					},
					{
						Name:     "cpu_options_threads_per_core",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("CpuOptions.ThreadsPerCore"),
					},
				},
			},
			{
				Name: "ebs_optimized",
				Type: schema.TypeBool,
			},
			{
				Name: "ena_support",
				Type: schema.TypeBool,
			},
			{
				Name: "enclave_options",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "enclave_options_enabled",
						Type:     schema.TypeBool,
						Resolver: schema.PathResolver("EnclaveOptions.Enabled"),
					},
				},
			},
			{
				Name: "hibernation_options",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "hibernation_options_configured",
						Type:     schema.TypeBool,
						Resolver: schema.PathResolver("HibernationOptions.Configured"),
					},
				},
			},
			{
				Name: "hypervisor",
				Type: schema.TypeString,
			},
			{
				Name: "iam_instance_profile",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "iam_instance_profile_arn",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("IamInstanceProfile.Arn"),
					},
					{
						Name:     "iam_instance_profile_id",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("IamInstanceProfile.Id"),
					},
				},
			},
			{
				Name: "image_id",
				Type: schema.TypeString,
			},
			{
				Name: "instance_id",
				Type: schema.TypeString,
			},
			{
				Name: "instance_lifecycle",
				Type: schema.TypeString,
			},
			{
				Name: "instance_type",
				Type: schema.TypeString,
			},
			{
				Name: "kernel_id",
				Type: schema.TypeString,
			},
			{
				Name: "key_name",
				Type: schema.TypeString,
			},
			{
				Name: "launch_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "metadata_options",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "metadata_options_http_endpoint",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("MetadataOptions.HttpEndpoint"),
					},
					{
						Name:     "metadata_options_http_put_response_hop_limit",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("MetadataOptions.HttpPutResponseHopLimit"),
					},
					{
						Name:     "metadata_options_http_tokens",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("MetadataOptions.HttpTokens"),
					},
					{
						Name:     "metadata_options_state",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("MetadataOptions.State"),
					},
				},
			},
			{
				Name: "monitoring",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "monitoring_state",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Monitoring.State"),
					},
				},
			},
			{
				Name: "outpost_arn",
				Type: schema.TypeString,
			},
			{
				Name: "placement",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "placement_affinity",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.Affinity"),
					},
					{
						Name:     "placement_availability_zone",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.AvailabilityZone"),
					},
					{
						Name:     "placement_group_name",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.GroupName"),
					},
					{
						Name:     "placement_host_id",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.HostId"),
					},
					{
						Name:     "placement_host_resource_group_arn",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.HostResourceGroupArn"),
					},
					{
						Name:     "placement_partition_number",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("Placement.PartitionNumber"),
					},
					{
						Name:     "placement_spread_domain",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.SpreadDomain"),
					},
					{
						Name:     "placement_tenancy",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Placement.Tenancy"),
					},
				},
			},
			{
				Name: "platform",
				Type: schema.TypeString,
			},
			{
				Name: "private_dns_name",
				Type: schema.TypeString,
			},
			{
				Name: "private_ip_address",
				Type: schema.TypeString,
			},
			{
				Name: "public_dns_name",
				Type: schema.TypeString,
			},
			{
				Name: "public_ip_address",
				Type: schema.TypeString,
			},
			{
				Name: "ramdisk_id",
				Type: schema.TypeString,
			},
			{
				Name: "root_device_name",
				Type: schema.TypeString,
			},
			{
				Name: "root_device_type",
				Type: schema.TypeString,
			},
			{
				Name: "source_dest_check",
				Type: schema.TypeBool,
			},
			{
				Name: "spot_instance_request_id",
				Type: schema.TypeString,
			},
			{
				Name: "sriov_net_support",
				Type: schema.TypeString,
			},
			{
				Name: "state",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "state_code",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("State.Code"),
					},
					{
						Name:     "state_name",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("State.Name"),
					},
				},
			},
			{
				Name: "state_reason",
				Type: schema.TypeEmbedded,
				Elem: []schema.Column{
					{
						Name:     "state_reason_code",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("StateReason.Code"),
					},
					{
						Name:     "state_reason_message",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("StateReason.Message"),
					},
				},
			},
			{
				Name: "state_transition_reason",
				Type: schema.TypeString,
			},
			{
				Name: "subnet_id",
				Type: schema.TypeString,
			},
			{
				Name: "tags",
				Type: schema.TypeJSON,
			},
			{
				Name: "virtualization_type",
				Type: schema.TypeString,
			},
			{
				Name: "vpc_id",
				Type: schema.TypeString,
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "aws_ec2_instance_block_device_mappings",
				Resolver: FetchEc2InstanceBlockDeviceMappings,
				Columns: []schema.Column{
					{
						Name: "device_name",
						Type: schema.TypeString,
					},
					{
						Name: "ebs",
						Type: schema.TypeEmbedded,
						Elem: []schema.Column{
							{
								Name:     "ebs_attach_time",
								Type:     schema.TypeTimestamp,
								Resolver: schema.PathResolver("Ebs.AttachTime"),
							},
							{
								Name:     "ebs_delete_on_termination",
								Type:     schema.TypeBool,
								Resolver: schema.PathResolver("Ebs.DeleteOnTermination"),
							},
							{
								Name:     "ebs_status",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Ebs.Status"),
							},
							{
								Name:     "ebs_volume_id",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Ebs.VolumeId"),
							},
						},
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
			{
				Name:     "aws_ec2_elastic_gpu_associations",
				Resolver: FetchEc2ElasticGpuAssociations,
				Columns: []schema.Column{
					{
						Name: "elastic_gpu_association_id",
						Type: schema.TypeString,
					},
					{
						Name: "elastic_gpu_association_state",
						Type: schema.TypeString,
					},
					{
						Name: "elastic_gpu_association_time",
						Type: schema.TypeString,
					},
					{
						Name: "elastic_gpu_id",
						Type: schema.TypeString,
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
			{
				Name:     "aws_ec2_elastic_inference_accelerator_associations",
				Resolver: FetchEc2ElasticInferenceAcceleratorAssociations,
				Columns: []schema.Column{
					{
						Name: "elastic_inference_accelerator_arn",
						Type: schema.TypeString,
					},
					{
						Name: "elastic_inference_accelerator_association_id",
						Type: schema.TypeString,
					},
					{
						Name: "elastic_inference_accelerator_association_state",
						Type: schema.TypeString,
					},
					{
						Name: "elastic_inference_accelerator_association_time",
						Type: schema.TypeTimestamp,
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
			{
				Name:     "aws_ec2_license_configurations",
				Resolver: FetchEc2LicenseConfigurations,
				Columns: []schema.Column{
					{
						Name: "license_configuration_arn",
						Type: schema.TypeString,
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
			{
				Name:     "aws_ec2_instance_network_interfaces",
				Resolver: FetchEc2InstanceNetworkInterfaces,
				Columns: []schema.Column{
					{
						Name: "association",
						Type: schema.TypeEmbedded,
						Elem: []schema.Column{
							{
								Name:     "association_carrier_ip",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Association.CarrierIp"),
							},
							{
								Name:     "association_ip_owner_id",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Association.IpOwnerId"),
							},
							{
								Name:     "association_public_dns_name",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Association.PublicDnsName"),
							},
							{
								Name:     "association_public_ip",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Association.PublicIp"),
							},
						},
					},
					{
						Name: "attachment",
						Type: schema.TypeEmbedded,
						Elem: []schema.Column{
							{
								Name:     "attachment_attach_time",
								Type:     schema.TypeTimestamp,
								Resolver: schema.PathResolver("Attachment.AttachTime"),
							},
							{
								Name:     "attachment_attachment_id",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Attachment.AttachmentId"),
							},
							{
								Name:     "attachment_delete_on_termination",
								Type:     schema.TypeBool,
								Resolver: schema.PathResolver("Attachment.DeleteOnTermination"),
							},
							{
								Name:     "attachment_device_index",
								Type:     schema.TypeInt,
								Resolver: schema.PathResolver("Attachment.DeviceIndex"),
							},
							{
								Name:     "attachment_network_card_index",
								Type:     schema.TypeInt,
								Resolver: schema.PathResolver("Attachment.NetworkCardIndex"),
							},
							{
								Name:     "attachment_status",
								Type:     schema.TypeString,
								Resolver: schema.PathResolver("Attachment.Status"),
							},
						},
					},
					{
						Name: "description",
						Type: schema.TypeString,
					},
					{
						Name: "interface_type",
						Type: schema.TypeString,
					},
					{
						Name: "mac_address",
						Type: schema.TypeString,
					},
					{
						Name: "network_interface_id",
						Type: schema.TypeString,
					},
					{
						Name: "owner_id",
						Type: schema.TypeString,
					},
					{
						Name: "private_dns_name",
						Type: schema.TypeString,
					},
					{
						Name: "private_ip_address",
						Type: schema.TypeString,
					},
					{
						Name: "source_dest_check",
						Type: schema.TypeBool,
					},
					{
						Name: "status",
						Type: schema.TypeString,
					},
					{
						Name: "subnet_id",
						Type: schema.TypeString,
					},
					{
						Name: "vpc_id",
						Type: schema.TypeString,
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
				Relations: []*schema.Table{
					{
						Name:     "aws_ec2_group_identifiers",
						Resolver: FetchEc2GroupIdentifiers,
						Columns: []schema.Column{
							{
								Name: "group_id",
								Type: schema.TypeString,
							},
							{
								Name: "group_name",
								Type: schema.TypeString,
							},
							{
								Name:     "instancenetworkinterface_id",
								Type:     schema.TypeUUID,
								Resolver: schema.ParentIdResolver,
							},
						},
					},
					{
						Name:     "aws_ec2_instance_ipv_6_addresses",
						Resolver: FetchEc2InstanceIpv6Addresses,
						Columns: []schema.Column{
							{
								Name: "ipv_6_address",
								Type: schema.TypeString,
							},
							{
								Name:     "instancenetworkinterface_id",
								Type:     schema.TypeUUID,
								Resolver: schema.ParentIdResolver,
							},
						},
					},
					{
						Name:     "aws_ec2_instance_private_ip_addresses",
						Resolver: FetchEc2InstancePrivateIPAddresses,
						Columns: []schema.Column{
							{
								Name: "association",
								Type: schema.TypeEmbedded,
								Elem: []schema.Column{
									{
										Name:     "association_carrier_ip",
										Type:     schema.TypeString,
										Resolver: schema.PathResolver("Association.CarrierIp"),
									},
									{
										Name:     "association_ip_owner_id",
										Type:     schema.TypeString,
										Resolver: schema.PathResolver("Association.IpOwnerId"),
									},
									{
										Name:     "association_public_dns_name",
										Type:     schema.TypeString,
										Resolver: schema.PathResolver("Association.PublicDnsName"),
									},
									{
										Name:     "association_public_ip",
										Type:     schema.TypeString,
										Resolver: schema.PathResolver("Association.PublicIp"),
									},
								},
							},
							{
								Name: "primary",
								Type: schema.TypeBool,
							},
							{
								Name: "private_dns_name",
								Type: schema.TypeString,
							},
							{
								Name: "private_ip_address",
								Type: schema.TypeString,
							},
							{
								Name:     "instancenetworkinterface_id",
								Type:     schema.TypeUUID,
								Resolver: schema.ParentIdResolver,
							},
						},
					},
				},
			},
			{
				Name:     "aws_ec2_product_codes",
				Resolver: FetchEc2ProductCodes,
				Columns: []schema.Column{
					{
						Name: "product_code_id",
						Type: schema.TypeString,
					},
					{
						Name: "product_code_type",
						Type: schema.TypeString,
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
			{
				Name:     "aws_ec2_group_identifiers",
				Resolver: FetchEc2GroupIdentifiers,
				Columns: []schema.Column{
					{
						Name: "group_id",
						Type: schema.TypeString,
					},
					{
						Name: "group_name",
						Type: schema.TypeString,
					},
					{
						Name:     "instance_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
				},
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================
func FetchEc2Instances(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2InstanceBlockDeviceMappings(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2ElasticGpuAssociations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2ElasticInferenceAcceleratorAssociations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2LicenseConfigurations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2InstanceNetworkInterfaces(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2GroupIdentifiers(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2InstanceIpv6Addresses(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2InstancePrivateIPAddresses(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
func FetchEc2ProductCodes(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	panic("no implemented")
}
