provider "ipam" {
  addresses = [
    {
      ip            = "1.1.1.1"
      prefix_length = "24"
      gateway       = "1.1.1.254"
    },
    {
      ip            = "1.1.1.2"
      prefix_length = "24"
      gateway       = "1.1.1.254"
    },
    {
      ip            = "1.1.1.3"
      prefix_length = "24"
      gateway       = "1.1.1.254"
    },
  ]
}
