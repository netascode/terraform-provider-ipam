package provider

// func TestAccIpamAllocate(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccIpamAllocateConfig_initial(),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("ipam_allocate.test", "addresses.host1.ip", "10.1.1.1"),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccIpamAllocateConfig_initial() string {
// 	return `
// 	resource "ipam_allocate" "test" {
// 		addresses = {
// 			"host1" = {}
// 		}
// 	}
// 	`
// }
