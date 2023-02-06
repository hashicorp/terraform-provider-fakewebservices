// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fws

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-fakewebservices/client"
)

// TODO: ADD DOCUMENTATION
// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: REMOVE THIS DEFAULT?
				DefaultFunc: schema.EnvDefaultFunc("FWS_HOSTNAME", client.DefaultHostname),
			},
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					// TODO: MAKE THIS WORK WITH DIFFERENT HOST
					// TODO: MAKE THIS WORK WITH OTHER CRED TYPES (THE REMOTE ONES)
					// TODO: PROBABLY REUSE terraform-svchost

					creds := *cliCredentials()
					return creds.Credentials["app.terraform.io"].Token, nil
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fakewebservices_server":        resourceFWSServer(),
			"fakewebservices_database":      resourceFWSDatabase(),
			"fakewebservices_load_balancer": resourceFWSLoadBalancer(),
			"fakewebservices_vpc":           resourceFWSVpc(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	hostname := d.Get("hostname").(string)
	token := d.Get("token").(string)
	return client.NewClient(hostname, token)
}
