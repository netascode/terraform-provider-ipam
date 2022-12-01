package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIpamAllocate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccIpamAllocateConfig_initial(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host1.ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host1.prefix_length", "22"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host1.gateway", "1.1.1.201"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host2.ip", "1.1.1.2"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host2.prefix_length", "22"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host2.gateway", "1.1.1.201"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host3.ip", "1.1.1.10"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host3.prefix_length", "23"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host3.gateway", "1.1.1.200"),
				),
			},
			{
				Config: providerConfig + testAccIpamAllocateConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host4.ip", "1.1.1.11"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host4.prefix_length", "24"),
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host4.gateway", "1.1.1.254"),
				),
			},
		},
	})
}

func testAccIpamAllocateConfig_initial() string {
	return `
	resource "ipam_allocate" "test" {
		pool = "POOL1"
		hosts = ["host1", "host2", "host3"]
	}
	`
}

func testAccIpamAllocateConfig_update() string {
	return `
	resource "ipam_allocate" "test" {
		pool = "POOL1"
		hosts = ["host1", "host2", "host3", "host4"]
	}
	`
}
