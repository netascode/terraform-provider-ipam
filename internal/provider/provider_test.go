package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ipam": func() (tfprotov6.ProviderServer, error) {
		provider := New("test")().(*provider)
		var addresses []providerAddress
		addresses = append(addresses, providerAddress{"10.1.1.1", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.2", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.3", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.4", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.5", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.6", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.7", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.8", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.9", "24", "10.1.1.254"})
		addresses = append(addresses, providerAddress{"10.1.1.10", "24", "10.1.1.254"})
		provider.addresses = addresses
		return tfsdk.NewProtocol6Server(provider), nil
	},
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
