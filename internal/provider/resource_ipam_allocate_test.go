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
					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host1.ip", "10.1.1.1"),
				),
			},
		},
	})
}

func testAccIpamAllocateConfig_initial() string {
	return `
	resource "ipam_allocate" "test" {
		addresses = {
			"host1" = {}
		}
	}
	`
}
