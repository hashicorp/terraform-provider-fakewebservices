---
page_title: "server Resource - terraform-provider-fakewebservices"
subcategory: ""
description: |-
  The order resource allows you to configure a FWS server.
---

# Resource `fakewebservices_server`

-> Visit the [Perform CRUD operations with Providers](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) Learn tutorial for an interactive getting started experience.

The order resource allows you to configure a HashiCups order.

## Example Usage

```terraform
resource "fakewebservices_server" "demoserver" {
  name = "demo-server"
}
```

## Argument Reference

- `name` - (Required) The name of the server
