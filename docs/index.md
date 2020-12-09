---
page_title: "Provider: Fake Web Services"
subcategory: ""
description: |-
  Terraform provider for interacting with Fake Web Services API.
---

# Fake Web Services Provider

The Fake Web Services Provider is used to provision demo infrastructure while onboarding to Terraform Cloud.

## Example Usage

```terraform
provider "fakewebservices" {
  host = "app.terraform.io"
}

resource "fakewebservices_server" "server" {
  name = "my-demo-server"
}

resource "fakewebservices_database" "database-prod" {
  name = "my-demo-database-prod"
}

resource "fakewebservices_database" "database-replica" {
  name = "my-demo-database-replica"
}
```

## Schema

### Optional

- **host** (String, Optional) FWS API address (defaults to `app.terraform.io`)
- **token** (String) Token used to authenticate into FWS API (defaults to `FWS_TOKEN` env var)
