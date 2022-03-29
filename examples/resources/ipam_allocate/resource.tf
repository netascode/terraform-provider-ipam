resource "ipam_allocate" "example" {
  addresses = {
    "host1" = {}
    "host2" = {}
  }
}
