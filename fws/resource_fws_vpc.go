package fws

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-fakewebservices/client"
)

func resourceFWSVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWSVpcCreate,
		Read:   resourceFWSVpcRead,
		Update: resourceFWSVpcUpdate,
		Delete: resourceFWSVpcDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the VPC.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cidr_block": {
				Description: "The range of IPv4 addresses for this VPC, in the form of a Classless Inter-Domain Routing (CIDR) block.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceFWSVpcCreate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	cb := d.Get("cidr_block").(string)

	options := VpcCreateOptions{
		Name:      client.String(name),
		CidrBlock: client.String(cb),
	}

	req, err := fwsClient.NewRequest("POST", "vpcs", &options)

	if err != nil {
		return err
	}

	vpc := &Vpc{}

	log.Printf("[DEBUG] Creating new vpc with name: %s", name)
	err = fwsClient.Do(req, vpc)
	if err != nil {
		return fmt.Errorf("Error creating vpc: %v", err)
	}

	d.SetId(vpc.ID)

	return resourceFWSVpcRead(d, meta)
}

func resourceFWSVpcRead(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest("GET", fmt.Sprintf("vpcs/%s", d.Id()), nil)

	if err != nil {
		return err
	}

	vpc := &Vpc{}

	log.Printf("[DEBUG] Reading vpc: %s", d.Id())
	err = fwsClient.Do(req, vpc)
	if err != nil {
		if err == client.ErrResourceNotFound {
			log.Printf("[DEBUG] vpc %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading configuration of vpc %s: %v", d.Id(), err)
	}

	// Update the config.
	d.Set("name", vpc.Name)

	return nil
}

func resourceFWSVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	name := d.Get("name").(string)
	cb := d.Get("cidr_block").(string)

	options := VpcUpdateOptions{
		Name:      client.String(name),
		CidrBlock: client.String(cb),
	}

	req, err := fwsClient.NewRequest(
		"PATCH",
		fmt.Sprintf("vpcs/%s", d.Id()),
		&options,
	)

	if err != nil {
		return err
	}

	vpc := &Vpc{}

	log.Printf("[DEBUG] Updating vpc: %s", d.Id())
	err = fwsClient.Do(req, vpc)
	if err != nil {
		return fmt.Errorf("Error updating vpc: %v", err)
	}

	return resourceFWSVpcRead(d, meta)
}

func resourceFWSVpcDelete(d *schema.ResourceData, meta interface{}) error {
	fwsClient := meta.(*client.Client)

	req, err := fwsClient.NewRequest(
		"DELETE",
		fmt.Sprintf("vpcs/%s", d.Id()),
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Destroying vpc: %s", d.Id())
	err = fwsClient.Do(req, nil)
	if err != nil {
		return fmt.Errorf("Error destroying vpc: %v", err)
	}

	return nil
}

type Vpc struct {
	ID        string `jsonapi:"primary,fake-resources-vpcs"`
	Name      string `jsonapi:"attr,name,omitempty"`
	CidrBlock string `jsonapi:"attr,cidr_block,omitempty"`
}

type VpcCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-vpcs"`

	// A name to identify the vpc.
	Name *string `jsonapi:"attr,name"`

	CidrBlock *string `jsonapi:"attr,cidr_block"`
}

type VpcUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,fake-resources-vpcs"`

	// A name to identify the vpc.
	Name *string `jsonapi:"attr,name"`

	CidrBlock *string `jsonapi:"attr,cidr_block"`
}
