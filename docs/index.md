---
layout: ""
page_title: "Provider: Fake Web Services"
description: |-
  Provision Terraform Cloud or Terraform Enterprise - with Terraform! Management of organizations, workspaces, teams, variables, run triggers, policy sets, and more. Maintained by the Terraform Cloud team at HashiCorp.
---

# "Fake Web Services" Provider

This provider provisions "resources" to a fictitious cloud provider, "Fake Web Services" - used in the [TFC Getting Started project](https://github.com/hashicorp/tfc-getting-started).

These resources are purely for demonstration and created in Terraform Cloud, scoped to your TFC user account.

## Installation & Usage

This provider isn't _really_ intended for any use beyond the example configuration, but you can absolutely use it outside the example if you like!

* Declare the provider in your configuration and `terraform init` will automatically fetch and install the provider for you from the [Terraform Registry](https://registry.terraform.io/).
* [Create a user or team API token in Terraform Cloud/Enterprise](https://www.terraform.io/docs/cloud/users-teams-organizations/api-tokens.html), and use the token in the provider configuration block.
* See the documentation for available resources and provision away!

Example:

```hcl
terraform {
  required_providers {
    fakewebservices = "~> 0.1"
  }
}

provider "fakewebservices" {
  token = var.provider_token
}

resource "fakewebservices_vpc" "primary_vpc" {
  name = "Primary VPC"
  cidr_block = "0.0.0.0/1"
}

resource "fakewebservices_server" "servers" {
  count = 2

  name = "Server ${count.index+1}"
  type = "t2.micro"
  vpc = fakewebservices_vpc.primary_vpc.name
}

resource "fakewebservices_load_balancer" "primary_lb" {
  name = "Primary Load Balancer"
  servers = fakewebservices_server.servers[*].name
}

resource "fakewebservices_database" "prod_db" {
  name = "Production DB"
  size = 256
}
```

## Schema

### Optional

- **hostname** (String)
- **token** (String)
