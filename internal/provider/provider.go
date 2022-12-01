package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider {
	return &ipamProvider{}
}

type ipamProvider struct{}

type providerAddress struct {
	ip           string
	prefixLength string
	gateway      string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	Pools []providerDataPool `tfsdk:"pools"`
}

type providerDataPool struct {
	Name         types.String              `tfsdk:"name"`
	PrefixLength types.Int64               `tfsdk:"prefix_length"`
	Gateway      types.String              `tfsdk:"gateway"`
	Ranges       []providerDataPoolRange   `tfsdk:"ranges"`
	Addresses    []providerDataPoolAddress `tfsdk:"addresses"`
}

type providerDataPoolRange struct {
	FromIP       types.String `tfsdk:"from_ip"`
	ToIP         types.String `tfsdk:"to_ip"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	Gateway      types.String `tfsdk:"gateway"`
}

type providerDataPoolAddress struct {
	IP           types.String `tfsdk:"ip"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	Gateway      types.String `tfsdk:"gateway"`
}

// Metadata returns the provider type name.
func (p *ipamProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ipam"
}

func (p *ipamProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"pools": {
				MarkdownDescription: "A list of managed IP pools.",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						MarkdownDescription: "IP pool name.",
						Type:                types.StringType,
						Required:            true,
					},
					"prefix_length": {
						MarkdownDescription: "Default prefix length.",
						Type:                types.Int64Type,
						Optional:            true,
					},
					"gateway": {
						MarkdownDescription: "Default gateway IP.",
						Type:                types.StringType,
						Optional:            true,
					},
					"ranges": {
						MarkdownDescription: "A list of IP ranges.",
						Optional:            true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"from_ip": {
								MarkdownDescription: "First IP.",
								Type:                types.StringType,
								Required:            true,
							},
							"to_ip": {
								MarkdownDescription: "Last IP.",
								Type:                types.StringType,
								Required:            true,
							},
							"prefix_length": {
								MarkdownDescription: "Prefix length.",
								Type:                types.Int64Type,
								Optional:            true,
							},
							"gateway": {
								MarkdownDescription: "Gateway IP.",
								Type:                types.StringType,
								Optional:            true,
							},
						}),
					},
					"addresses": {
						MarkdownDescription: "A list of IP addresses.",
						Optional:            true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"ip": {
								MarkdownDescription: "IP address.",
								Type:                types.StringType,
								Required:            true,
							},
							"prefix_length": {
								MarkdownDescription: "Prefix length.",
								Type:                types.Int64Type,
								Optional:            true,
							},
							"gateway": {
								MarkdownDescription: "Gateway IP.",
								Type:                types.StringType,
								Optional:            true,
							},
						}),
					},
				}),
			},
		},
	}, nil
}

func (p *ipamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for p := range config.Pools {
		globalPrefixLength := false
		if !config.Pools[p].PrefixLength.IsNull() {
			globalPrefixLength = true
			if err := ValidatePrefixLength(config.Pools[p].PrefixLength.ValueInt64()); err {
				resp.Diagnostics.AddError(
					"Invalid 'prefix_length' configured.",
					fmt.Sprintf("'prefix_length' must be a number between 0 and 128."),
				)
				return
			}
		}
		globalGateway := false
		if !config.Pools[p].Gateway.IsNull() {
			globalGateway = true
			if err := ValidateIPAddress(config.Pools[p].Gateway.ValueString()); err {
				resp.Diagnostics.AddError(
					"Invalid 'gateway' configured.",
					fmt.Sprintf("'gateway' is not a valid IP address."),
				)
				return
			}
		}
		for r := range config.Pools[p].Ranges {
			if !config.Pools[p].Ranges[r].PrefixLength.IsNull() {
				if err := ValidatePrefixLength(config.Pools[p].Ranges[r].PrefixLength.ValueInt64()); err {
					resp.Diagnostics.AddError(
						"Invalid 'prefix_length' configured.",
						fmt.Sprintf("'prefix_length' must be a number between 0 and 128."),
					)
					return
				}
			}
			if !globalPrefixLength && config.Pools[p].Ranges[r].PrefixLength.IsNull() {
				resp.Diagnostics.AddError(
					"Range without 'prefix_length' configured.",
					fmt.Sprintf("Range '%s-%s' has no 'prefix_length' configured.", config.Pools[p].Ranges[r].FromIP.ValueString(), config.Pools[p].Ranges[r].ToIP.ValueString()),
				)
				return
			}
			if !config.Pools[p].Ranges[r].Gateway.IsNull() {
				if err := ValidateIPAddress(config.Pools[p].Ranges[r].Gateway.ValueString()); err {
					resp.Diagnostics.AddError(
						"Invalid 'gateway' configured.",
						fmt.Sprintf("'gateway' is not a valid IP address."),
					)
					return
				}
			}
			if !globalGateway && config.Pools[p].Ranges[r].Gateway.IsNull() {
				resp.Diagnostics.AddError(
					"Range without 'gateway' configured.",
					fmt.Sprintf("Range '%s-%s' has no 'gateway' configured.", config.Pools[p].Ranges[r].FromIP.ValueString(), config.Pools[p].Ranges[r].ToIP.ValueString()),
				)
				return
			}
			if err := ValidateIPAddress(config.Pools[p].Ranges[r].FromIP.ValueString()); err {
				resp.Diagnostics.AddError(
					"Invalid 'from_ip' configured.",
					fmt.Sprintf("IP '%s' is not a valid address.", config.Pools[p].Ranges[r].FromIP.ValueString()),
				)
				return
			}
			if err := ValidateIPAddress(config.Pools[p].Ranges[r].ToIP.ValueString()); err {
				resp.Diagnostics.AddError(
					"Invalid 'to_ip' configured.",
					fmt.Sprintf("IP '%s' is not a valid address.", config.Pools[p].Ranges[r].ToIP.ValueString()),
				)
				return
			}
			if err := ValidateIPRange(config.Pools[p].Ranges[r].FromIP.ValueString(), config.Pools[p].Ranges[r].ToIP.ValueString()); err {
				resp.Diagnostics.AddError(
					"Invalid range configured.",
					fmt.Sprintf("Range '%s-%s', 'from_ip' must be smaller than 'to_ip'.", config.Pools[p].Ranges[r].FromIP.ValueString(), config.Pools[p].Ranges[r].ToIP.ValueString()),
				)
				return
			}
		}
		for a := range config.Pools[p].Addresses {
			if !config.Pools[p].Addresses[a].PrefixLength.IsNull() {
				if err := ValidatePrefixLength(config.Pools[p].Addresses[a].PrefixLength.ValueInt64()); err {
					resp.Diagnostics.AddError(
						"Invalid 'prefix_length' configured.",
						fmt.Sprintf("'prefix_length' must be a number between 0 and 128."),
					)
					return
				}
			}
			if !globalPrefixLength && config.Pools[p].Addresses[a].PrefixLength.IsNull() {
				resp.Diagnostics.AddError(
					"Address without 'prefix_length' configured.",
					fmt.Sprintf("IP '%s' has no 'prefix_length' configured.", config.Pools[p].Addresses[a].IP.ValueString()),
				)
				return
			}
			if !config.Pools[p].Addresses[a].Gateway.IsNull() {
				if err := ValidateIPAddress(config.Pools[p].Addresses[a].Gateway.ValueString()); err {
					resp.Diagnostics.AddError(
						"Invalid 'gateway' configured.",
						fmt.Sprintf("'gateway' is not a valid IP address."),
					)
					return
				}
			}
			if !globalGateway && config.Pools[p].Addresses[a].Gateway.IsNull() {
				resp.Diagnostics.AddError(
					"Address without 'gateway' configured.",
					fmt.Sprintf("IP '%s' has no 'gateway' configured.", config.Pools[p].Addresses[a].IP.ValueString()),
				)
				return
			}
			if err := ValidateIPAddress(config.Pools[p].Addresses[a].IP.ValueString()); err {
				resp.Diagnostics.AddError(
					"Invalid 'ip' configured.",
					fmt.Sprintf("IP '%s' is not a valid address.", config.Pools[p].Addresses[a].IP.ValueString()),
				)
				return
			}
		}
	}

	resp.DataSourceData = config.Pools
	resp.ResourceData = config.Pools
}

func (p *ipamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIpamAllocateResource,
	}
}

func (p *ipamProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
