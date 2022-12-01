provider "ipam" {
  pools = [
    {
      name          = "POOL1"
      prefix_length = 24
      gateway       = "1.1.1.254"
      ranges = [
        {
          from_ip = "1.1.1.1"
          to_ip   = "1.1.1.10"
        }
      ]
      addresses = [
        {
          ip = "1.1.1.20"
        },
        {
          ip = "1.1.1.30"
        },
      ]
    }
  ]
}
