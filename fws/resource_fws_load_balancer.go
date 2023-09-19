// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fws

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-fakewebservices/client"
)

func resourceFWSLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWSLoadBalancerCreate,
		Read:   resourceFWSLoadBalancerRead,
		Update: resourceFWSLoadBalancerUpdate,
		Delete: resourceFWSLoadBalancerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the load balancer.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"servers": {
				Description: "A list of server names to attach to the load balancer.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
		},
	}
}

func resourceFWSLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	servers := client.ExpandStringSet(d.Get("servers").(*schema.Set))

	options := LoadBalancerCreateOptions{
		Name:    client.String(name),
		Servers: &servers,
	}

	req, err := fwsClient.NewRequest("POST", "load_balancers", &options)

	if err != nil {
		return err
	}

	lb := &LoadBalancer{}

	log.Printf("[DEBUG] Creating new load_balancer with name: %s", name)
	err = fwsClient.Do(req, lb)
	if err != nil {
		return fmt.Errorf("Error creating load_balancer: %v", err)
	}

	d.SetId(lb.ID)

	return resourceFWSLoadBalancerRead(d, meta)
}

func resourceFWSLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest("GET", fmt.Sprintf("load_balancers/%s", d.Id()), nil)

	if err != nil {
		return err
	}

	lb := &LoadBalancer{}

	log.Printf("[DEBUG] Reading load_balancer: %s", d.Id())
	err = fwsClient.Do(req, lb)
	if err != nil {
		if err == client.ErrResourceNotFound {
			log.Printf("[DEBUG] load_balancer %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading configuration of load_balancer %s: %v", d.Id(), err)
	}

	// Update the config.
	d.Set("name", lb.Name)

	return nil
}

func resourceFWSLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	servers := client.ExpandStringSet(d.Get("servers").(*schema.Set))

	options := LoadBalancerUpdateOptions{
		Name:    client.String(name),
		Servers: &servers,
	}

	req, err := fwsClient.NewRequest(
		"PATCH",
		fmt.Sprintf("load_balancers/%s", d.Id()),
		&options,
	)

	if err != nil {
		return err
	}

	lb := &LoadBalancer{}

	log.Printf("[DEBUG] Updating load_balancer: %s", d.Id())
	err = fwsClient.Do(req, lb)
	if err != nil {
		return fmt.Errorf("Error updating load_balancer: %v", err)
	}

	return resourceFWSLoadBalancerRead(d, meta)
}

func resourceFWSLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest(
		"DELETE",
		fmt.Sprintf("load_balancers/%s", d.Id()),
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Destroying load_balancer: %s", d.Id())
	err = fwsClient.Do(req, nil)
	if err != nil {
		return fmt.Errorf("Error destroying load_balancer: %v", err)
	}

	return nil
}

type LoadBalancer struct {
	ID      string   `jsonapi:"primary,fake-resources-load-balancers"`
	Name    string   `jsonapi:"attr,name,omitempty"`
	Servers []string `jsonapi:"attr,servers,omitempty"`
}

type LoadBalancerCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-load-balancers"`

	Name    *string   `jsonapi:"attr,name"`
	Servers *[]string `jsonapi:"attr,servers"`
}

type LoadBalancerUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-load-balancers"`

	Name    *string   `jsonapi:"attr,name"`
	Servers *[]string `jsonapi:"attr,servers"`
}
