package provider

import (
	"context"

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
	Addresses []providerDataAddress `tfsdk:"addresses"`
}

type providerDataAddress struct {
	IP           types.String `tfsdk:"ip"`
	PrefixLength types.String `tfsdk:"prefix_length"`
	Gateway      types.String `tfsdk:"gateway"`
}

// Metadata returns the provider type name.
func (p *ipamProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ipam"
}

func (p *ipamProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"addresses": {
				MarkdownDescription: "A list of managed addresses.",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"ip": {
						MarkdownDescription: "IP address.",
						Type:                types.StringType,
						Required:            true,
					},
					"prefix_length": {
						MarkdownDescription: "Prefix length.",
						Type:                types.StringType,
						Required:            true,
					},
					"gateway": {
						MarkdownDescription: "Gateway IP.",
						Type:                types.StringType,
						Required:            true,
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

	var addresses []providerAddress
	for _, a := range config.Addresses {
		addresses = append(addresses, providerAddress{ip: a.IP.ValueString(), prefixLength: a.PrefixLength.ValueString(), gateway: a.Gateway.ValueString()})
	}

	resp.DataSourceData = addresses
	resp.ResourceData = addresses
}

func (p *ipamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIpamAllocateResource,
	}
}

func (p *ipamProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
