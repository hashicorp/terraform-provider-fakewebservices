package fws

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-fakewebservices/client"
)

func resourceFWSDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWSDatabaseCreate,
		Read:   resourceFWSDatabaseRead,
		Update: resourceFWSDatabaseUpdate,
		Delete: resourceFWSDatabaseDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the database.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The allocated size of the database in gigabytes.",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func resourceFWSDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	size := d.Get("size").(int)

	options := DatabaseCreateOptions{
		Name: client.String(name),
		Size: client.Int(size),
	}

	req, err := fwsClient.NewRequest("POST", "databases", &options)

	if err != nil {
		return err
	}

	database := &Database{}

	log.Printf("[DEBUG] Creating new database with name: %s", name)
	err = fwsClient.Do(req, database)
	if err != nil {
		return fmt.Errorf("Error creating database: %v", err)
	}

	d.SetId(database.ID)

	return resourceFWSDatabaseRead(d, meta)
}

func resourceFWSDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest("GET", fmt.Sprintf("databases/%s", d.Id()), nil)

	if err != nil {
		return err
	}

	database := &Database{}

	log.Printf("[DEBUG] Reading database: %s", d.Id())
	err = fwsClient.Do(req, database)
	if err != nil {
		if err == client.ErrResourceNotFound {
			log.Printf("[DEBUG] database %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading configuration of database %s: %v", d.Id(), err)
	}

	// Update the config.
	d.Set("name", database.Name)

	return nil
}

func resourceFWSDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	size := d.Get("size").(int)

	options := DatabaseUpdateOptions{
		Name: client.String(name),
		Size: client.Int(size),
	}

	req, err := fwsClient.NewRequest(
		"PATCH",
		fmt.Sprintf("databases/%s", d.Id()),
		&options,
	)

	if err != nil {
		return err
	}

	database := &Database{}

	log.Printf("[DEBUG] Updating database: %s", d.Id())
	err = fwsClient.Do(req, database)
	if err != nil {
		return fmt.Errorf("Error updating database: %v", err)
	}

	return resourceFWSDatabaseRead(d, meta)
}

func resourceFWSDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest(
		"DELETE",
		fmt.Sprintf("databases/%s", d.Id()),
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Destroying database: %s", d.Id())
	err = fwsClient.Do(req, nil)
	if err != nil {
		return fmt.Errorf("Error destroying database: %v", err)
	}

	return nil
}

type Database struct {
	ID   string `jsonapi:"primary,fake-resources-databases"`
	Name string `jsonapi:"attr,name,omitempty"`
	Size int    `jsonapi:"attr,size,omitempty"`
}

type DatabaseCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-databases"`

	// A name to identify the database.
	Name *string `jsonapi:"attr,name"`

	Size *int `jsonapi:"attr,size"`
}

type DatabaseUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-databases"`

	// A name to identify the database.
	Name *string `jsonapi:"attr,name"`

	Size *int `jsonapi:"attr,size"`
}
