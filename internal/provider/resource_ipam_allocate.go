package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = (*ipamAllocateResource)(nil)

func NewIpamAllocateResource() resource.Resource {
	return &ipamAllocateResource{}
}

type ipamAllocateResource struct {
	pools []providerDataPool
}

func (r *ipamAllocateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_allocate"
}

func (r *ipamAllocateResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Allocate one IP from a pool per unique host ID. A single resource must be used per pool.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Random internal ID.",
				Type:        types.StringType,
				Computed:    true,
			},
			"pool": {
				Description: "Pool name.",
				Type:        types.StringType,
				Required:    true,
			},
			"hosts": {
				Description: "List of unique host IDs.",
				Type:        types.ListType{ElemType: types.StringType},
				Required:    true,
			},
			"addresses": {
				MarkdownDescription: "A map of host IDs and its assigned addresses.",
				Computed:            true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"ip": {
						MarkdownDescription: "IP address.",
						Type:                types.StringType,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"prefix_length": {
						MarkdownDescription: "Prefix length.",
						Type:                types.Int64Type,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"gateway": {
						MarkdownDescription: "Gateway IP.",
						Type:                types.StringType,
						Computed:            true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
				}),
			},
		},
	}, nil
}

type Allocate struct {
	Id        types.String `tfsdk:"id"`
	Pool      types.String `tfsdk:"pool"`
	Hosts     types.List   `tfsdk:"hosts"`
	Addresses types.Map    `tfsdk:"addresses"`
}

var AddressType = map[string]attr.Type{"gateway": types.StringType, "ip": types.StringType, "prefix_length": types.Int64Type}

func (r *ipamAllocateResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.pools = req.ProviderData.([]providerDataPool)
}

func (r *ipamAllocateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state Allocate

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Beginning Create"))

	var pool *providerDataPool

	for i := range r.pools {
		if r.pools[i].Name.ValueString() == plan.Pool.ValueString() {
			pool = &r.pools[i]
		}
	}
	if pool == nil {
		resp.Diagnostics.AddError("Pool not found", fmt.Sprintf("Pool '%s' not found.", plan.Pool.ValueString()))
		return
	}

	poolAddresses := GetAddressesFromPool(pool)

	hosts := plan.Hosts.Elements()
	addresses := plan.Addresses.Elements()
	if addresses == nil {
		addresses = make(map[string]attr.Value)
	}

	if len(hosts) > len(poolAddresses) {
		resp.Diagnostics.AddError("Not enough IPs in pool", fmt.Sprintf("Pool '%s' does not have enough IP addresses.", plan.Pool.ValueString()))
		return
	}

	index := 0
	for h := range hosts {
		host := hosts[h].(types.String)
		ip := poolAddresses[index].IP
		prefixLength := poolAddresses[index].PrefixLength
		gateway := poolAddresses[index].Gateway
		v := map[string]attr.Value{"gateway": gateway, "ip": ip, "prefix_length": prefixLength}
		o := types.ObjectValueMust(AddressType, v)
		addresses[host.ValueString()] = o
		tflog.Debug(ctx, fmt.Sprintf("Allocate IP to %s: %v", host, ip.ValueString()))
		index += 1
	}

	state.Pool = plan.Pool
	state.Hosts = plan.Hosts
	state.Addresses = types.MapValueMust(plan.Addresses.ElementType(ctx), addresses)

	rand.Seed(time.Now().UnixNano())
	state.Id = types.StringValue(fmt.Sprint(rand.Int63()))

	tflog.Debug(ctx, fmt.Sprintf("Create finished successfully"))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ipamAllocateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Allocate

	// Read state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Beginning Read"))

	tflog.Debug(ctx, fmt.Sprintf("Read finished successfully"))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ipamAllocateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state Allocate

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Beginning Update"))

	var pool *providerDataPool

	for i := range r.pools {
		if r.pools[i].Name.ValueString() == plan.Pool.ValueString() {
			pool = &r.pools[i]
		}
	}
	if pool == nil {
		resp.Diagnostics.AddError("Pool not found", fmt.Sprintf("Pool '%s' not found.", plan.Pool.ValueString()))
		return
	}

	poolAddresses := GetAddressesFromPool(pool)

	hosts := plan.Hosts.Elements()
	addresses := plan.Addresses.Elements()
	if addresses == nil {
		addresses = make(map[string]attr.Value)
	}

	if len(hosts) > len(poolAddresses) {
		resp.Diagnostics.AddError("Not enough IPs in pool", fmt.Sprintf("Pool '%s' does not have enough IP addresses.", plan.Pool.ValueString()))
		return
	}

	for h := range hosts {
		host := hosts[h].(types.String)
		// check if an address is already assigned
		found := false
		var inUseAddresses []string
		for h, a := range addresses {
			address := a.(types.Object)
			if h == host.ValueString() {
				found = true
			}
			inUseAddresses = append(inUseAddresses, address.Attributes()["ip"].(types.String).ValueString())
		}
		if found {
			continue
		}
		// find next free IP
		for pa := range poolAddresses {
			inUse := false
			for _, inUseAddress := range inUseAddresses {
				if poolAddresses[pa].IP.ValueString() == inUseAddress {
					inUse = true
					break
				}
			}
			if inUse {
				continue
			}
			ip := poolAddresses[pa].IP
			prefixLength := poolAddresses[pa].PrefixLength
			gateway := poolAddresses[pa].Gateway
			v := map[string]attr.Value{"gateway": gateway, "ip": ip, "prefix_length": prefixLength}
			o := types.ObjectValueMust(AddressType, v)
			addresses[host.ValueString()] = o
			tflog.Debug(ctx, fmt.Sprintf("Allocate IP to %s: %v", host, ip.ValueString()))
			break
		}
	}

	state.Id = types.StringValue(plan.Id.ValueString())
	state.Pool = plan.Pool
	state.Hosts = plan.Hosts
	state.Addresses = types.MapValueMust(plan.Addresses.ElementType(ctx), addresses)

	tflog.Debug(ctx, fmt.Sprintf("Update finished successfully"))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ipamAllocateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Allocate

	// Read state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Beginning Delete"))

	tflog.Debug(ctx, fmt.Sprintf("Delete finished successfully"))

	resp.State.RemoveResource(ctx)
}
