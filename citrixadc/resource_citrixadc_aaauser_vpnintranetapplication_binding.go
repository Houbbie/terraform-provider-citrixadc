package citrixadc

import (
	"github.com/citrix/adc-nitro-go/resource/config/aaa"
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"log"
	"strings"
)

func resourceCitrixAdcAaauser_vpnintranetapplication_binding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createAaauser_vpnintranetapplication_bindingFunc,
		Read:          readAaauser_vpnintranetapplication_bindingFunc,
		Delete:        deleteAaauser_vpnintranetapplication_bindingFunc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"intranetapplication": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gotopriorityexpression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func createAaauser_vpnintranetapplication_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In createAaauser_vpnintranetapplication_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	username := d.Get("username").(string)
	intranetapplication := d.Get("intranetapplication").(string)
	bindingId := fmt.Sprintf("%s,%s", username, intranetapplication)
	aaauser_vpnintranetapplication_binding := aaa.Aaauservpnintranetapplicationbinding{
		Gotopriorityexpression: d.Get("gotopriorityexpression").(string),
		Intranetapplication:    d.Get("intranetapplication").(string),
		Username:               d.Get("username").(string),
	}

	err := client.UpdateUnnamedResource(service.Aaauser_vpnintranetapplication_binding.Type(), &aaauser_vpnintranetapplication_binding)
	if err != nil {
		return err
	}

	d.SetId(bindingId)

	err = readAaauser_vpnintranetapplication_bindingFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this aaauser_vpnintranetapplication_binding but we can't read it ?? %s", bindingId)
		return nil
	}
	return nil
}

func readAaauser_vpnintranetapplication_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] citrixadc-provider:  In readAaauser_vpnintranetapplication_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	username := idSlice[0]
	intranetapplication := idSlice[1]

	log.Printf("[DEBUG] citrixadc-provider: Reading aaauser_vpnintranetapplication_binding state %s", bindingId)

	findParams := service.FindParams{
		ResourceType:             "aaauser_vpnintranetapplication_binding",
		ResourceName:             username,
		ResourceMissingErrorCode: 258,
	}
	dataArr, err := client.FindResourceArrayWithParams(findParams)

	// Unexpected error
	if err != nil {
		log.Printf("[DEBUG] citrixadc-provider: Error during FindResourceArrayWithParams %s", err.Error())
		return err
	}

	// Resource is missing
	if len(dataArr) == 0 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams returned empty array")
		log.Printf("[WARN] citrixadc-provider: Clearing aaauser_vpnintranetapplication_binding state %s", bindingId)
		d.SetId("")
		return nil
	}

	// Iterate through results to find the one with the right id
	foundIndex := -1
	for i, v := range dataArr {
		if v["intranetapplication"].(string) == intranetapplication {
			foundIndex = i
			break
		}
	}

	// Resource is missing
	if foundIndex == -1 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams intranetapplication not found in array")
		log.Printf("[WARN] citrixadc-provider: Clearing aaauser_vpnintranetapplication_binding state %s", bindingId)
		d.SetId("")
		return nil
	}
	// Fallthrough

	data := dataArr[foundIndex]

	d.Set("gotopriorityexpression", data["gotopriorityexpression"])
	d.Set("intranetapplication", data["intranetapplication"])
	d.Set("username", data["username"])

	return nil

}

func deleteAaauser_vpnintranetapplication_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In deleteAaauser_vpnintranetapplication_bindingFunc")
	client := meta.(*NetScalerNitroClient).client

	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	name := idSlice[0]
	intranetapplication := idSlice[1]

	args := make([]string, 0)
	args = append(args, fmt.Sprintf("intranetapplication:%s", intranetapplication))

	err := client.DeleteResourceWithArgs(service.Aaauser_vpnintranetapplication_binding.Type(), name, args)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
