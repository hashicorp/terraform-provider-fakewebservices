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
		},
	}
}

func resourceFWSServerCreate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	serverType := d.Get("type").(string)

	options := ServerCreateOptions{
		Name:       client.String(name),
		ServerType: client.String(serverType),
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

	return nil
}

func resourceFWSServerUpdate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	serverType := d.Get("type").(string)

	options := ServerUpdateOptions{
		Name:       client.String(name),
		ServerType: client.String(serverType),
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
	ID         string `jsonapi:"primary,fake-resources-servers"`
	Name       string `jsonapi:"attr,name,omitempty"`
	ServerType string `jsonapi:"attr,server_type,omitempty"`
}

type ServerCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-servers"`

	// A name to identify the server.
	Name       *string `jsonapi:"attr,name"`
	ServerType *string `jsonapi:"attr,server_type"`
}

type ServerUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-servers"`

	// A name to identify the server.
	Name       *string `jsonapi:"attr,name"`
	ServerType *string `jsonapi:"attr,server_type"`
}
