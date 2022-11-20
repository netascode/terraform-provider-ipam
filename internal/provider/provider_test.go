package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "ipam" {
  addresses = [
    {
    	ip            = "10.1.1.1"
    	prefix_length = "24"
    	 gateway       = "10.1.1.254"
    },
    {
		ip            = "10.1.1.2"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.3"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.4"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.5"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.6"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.7"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.8"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.9"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
	{
		ip            = "10.1.1.10"
		prefix_length = "24"
		gateway       = "10.1.1.254"
	},
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
