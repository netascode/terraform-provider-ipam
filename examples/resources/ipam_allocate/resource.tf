resource "ipam_allocate" "example" {
  pool = "POOL1"
  hosts = {
    "host1" = {}
    "host2" = {}
  }
}

output "hosts" {
  value = ipam_allocate.example.hosts
}

/* 
hosts = tomap({
  "host1" = {
    "gateway" = "1.1.1.254"
    "ip" = "1.1.1.1"
    "prefix_length" = 24
  }
  "host2" = {
    "gateway" = "1.1.1.254"
    "ip" = "1.1.1.2"
    "prefix_length" = 24
  }
})
*/
