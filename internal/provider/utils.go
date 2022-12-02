package provider

import (
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetAddressesFromPool(pool *providerDataPool) []providerDataPoolAddress {
	addresses := make([]providerDataPoolAddress, 0)
	globalPrefixLength := pool.PrefixLength.ValueInt64()
	globalGateway := pool.Gateway.ValueString()
	for r := range pool.Ranges {
		fromIp, _ := netip.ParseAddr(pool.Ranges[r].FromIP.ValueString())
		toIp, _ := netip.ParseAddr(pool.Ranges[r].ToIP.ValueString())
		var prefixLength int64
		var gateway string
		if pool.Ranges[r].PrefixLength.IsNull() {
			prefixLength = globalPrefixLength
		} else {
			prefixLength = pool.Ranges[r].PrefixLength.ValueInt64()
		}
		if pool.Ranges[r].Gateway.IsNull() {
			gateway = globalGateway
		} else {
			gateway = pool.Ranges[r].Gateway.ValueString()
		}
		ip := fromIp
		for ip.Less(toIp) || ip == toIp {
			addresses = append(addresses, providerDataPoolAddress{IP: types.StringValue(ip.String()), PrefixLength: types.Int64Value(prefixLength), Gateway: types.StringValue(gateway)})
			ip = ip.Next()
		}
	}
	for a := range pool.Addresses {
		var prefixLength int64
		var gateway string
		if pool.Addresses[a].PrefixLength.IsNull() {
			prefixLength = globalPrefixLength
		} else {
			prefixLength = pool.Addresses[a].PrefixLength.ValueInt64()
		}
		if pool.Addresses[a].Gateway.IsNull() {
			gateway = globalGateway
		} else {
			gateway = pool.Addresses[a].Gateway.ValueString()
		}
		ip := pool.Addresses[a].IP
		addresses = append(addresses, providerDataPoolAddress{IP: ip, PrefixLength: types.Int64Value(prefixLength), Gateway: types.StringValue(gateway)})
	}
	return addresses
}

func ValidateIPAddress(ip string) bool {
	if _, err := netip.ParseAddr(ip); err != nil {
		return true
	}
	return false
}

func ValidatePrefixLength(prefixLength int64) bool {
	if prefixLength < 0 || prefixLength > 128 {
		return true
	}
	return false
}

func ValidateIPRange(fromIp, toIp string) bool {
	f := netip.MustParseAddr(fromIp)
	t := netip.MustParseAddr(toIp)
	if !f.Less(t) {
		return true
	}
	return false
}
