package fws

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-fakewebservices/client"
)

func resourceFWSServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWSServerCreate,
		Read:   resourceFWSServerRead,
		Update: resourceFWSServerUpdate,
		Delete: resourceFWSServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the server.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The server type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vpc": {
				Description: "The name of the VPC to deploy this server in.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceFWSServerCreate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	serverType := d.Get("type").(string)
	vpc := d.Get("vpc").(string)

	options := ServerCreateOptions{
		Name: client.String(name),
		Type: client.String(serverType),
		VPC:  client.String(vpc),
	}

	req, err := fwsClient.NewRequest("POST", "servers", &options)

	if err != nil {
		return err
	}

	server := &Server{}

	log.Printf("[DEBUG] Creating new server with name: %s", name)
	err = fwsClient.Do(req, server)
	if err != nil {
		return fmt.Errorf("Error creating server: %v", err)
	}

	d.SetId(server.ID)

	return resourceFWSServerRead(d, meta)
}

func resourceFWSServerRead(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest("GET", fmt.Sprintf("servers/%s", d.Id()), nil)

	if err != nil {
		return err
	}

	server := &Server{}

	log.Printf("[DEBUG] Reading server: %s", d.Id())
	err = fwsClient.Do(req, server)
	if err != nil {
		if err == client.ErrResourceNotFound {
			log.Printf("[DEBUG] server %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading configuration of server %s: %v", d.Id(), err)
	}

	// Update the config.
	d.Set("name", server.Name)
	d.Set("type", server.Type)
	d.Set("vpc", server.VPC)

	return nil
}

func resourceFWSServerUpdate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	serverType := d.Get("type").(string)
	vpc := d.Get("vpc").(string)

	options := ServerUpdateOptions{
		Name: client.String(name),
		Type: client.String(serverType),
		VPC:  client.String(vpc),
	}

	req, err := fwsClient.NewRequest(
		"PATCH",
		fmt.Sprintf("servers/%s", d.Id()),
		&options,
	)

	if err != nil {
		return err
	}

	server := &Server{}

	log.Printf("[DEBUG] Updating server: %s", d.Id())
	err = fwsClient.Do(req, server)
	if err != nil {
		return fmt.Errorf("Error updating server: %v", err)
	}

	return resourceFWSServerRead(d, meta)
}

func resourceFWSServerDelete(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest(
		"DELETE",
		fmt.Sprintf("servers/%s", d.Id()),
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Destroying server: %s", d.Id())
	err = fwsClient.Do(req, nil)
	if err != nil {
		return fmt.Errorf("Error destroying server: %v", err)
	}

	return nil
}

type Server struct {
	ID string `jsonapi:"primary,fake-resources-servers"`

	Name string `jsonapi:"attr,name,omitempty"`
	Type string `jsonapi:"attr,server-type,omitempty"`
	VPC  string `jsonapi:"attr,vpc,omitempty"`
}

type ServerCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-servers"`

	Name *string `jsonapi:"attr,name"`
	Type *string `jsonapi:"attr,server-type"`
	VPC  *string `jsonapi:"attr,vpc"`
}

type ServerUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-servers"`

	Name *string `jsonapi:"attr,name"`
	Type *string `jsonapi:"attr,server-type"`
	VPC  *string `jsonapi:"attr,vpc"`
}
