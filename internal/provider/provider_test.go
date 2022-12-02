package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "ipam" {
	pools = [
		{
			name = "POOL1"
			prefix_length = 24
			gateway = "1.1.1.254"
			ranges = [
				{
					from_ip = "1.1.1.1"
					to_ip = "1.1.1.2"
					prefix_length = 22
					gateway = "1.1.1.201"
				}
			]
			addresses = [
				{
					ip            = "1.1.1.10"
					prefix_length = 23
					gateway       = "1.1.1.200"
				},
				{
					ip            = "1.1.1.11"
				},
			]
		}
	]
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"ipam": providerserver.NewProtocol6WithError(New()),
	}
)
