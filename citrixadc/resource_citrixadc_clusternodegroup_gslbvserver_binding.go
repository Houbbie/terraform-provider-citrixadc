package citrixadc

import (
	"github.com/citrix/adc-nitro-go/resource/config/cluster"
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"log"
	"strings"
)

func resourceCitrixAdcClusternodegroup_gslbvserver_binding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createClusternodegroup_gslbvserver_bindingFunc,
		Read:          readClusternodegroup_gslbvserver_bindingFunc,
		Delete:        deleteClusternodegroup_gslbvserver_bindingFunc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vserver": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createClusternodegroup_gslbvserver_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In createClusternodegroup_gslbvserver_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	name := d.Get("name")
	vserver := d.Get("vserver")
	bindingId := fmt.Sprintf("%s,%s", name, vserver)
	clusternodegroup_gslbvserver_binding := cluster.Clusternodegroupgslbvserverbinding{
		Name:    d.Get("name").(string),
		Vserver: d.Get("vserver").(string),
	}

	err := client.UpdateUnnamedResource(service.Clusternodegroup_gslbvserver_binding.Type(), &clusternodegroup_gslbvserver_binding)
	if err != nil {
		return err
	}

	d.SetId(bindingId)

	err = readClusternodegroup_gslbvserver_bindingFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this clusternodegroup_gslbvserver_binding but we can't read it ?? %s", bindingId)
		return nil
	}
	return nil
}

func readClusternodegroup_gslbvserver_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] citrixadc-provider:  In readClusternodegroup_gslbvserver_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	name := idSlice[0]
	vserver := idSlice[1]

	log.Printf("[DEBUG] citrixadc-provider: Reading clusternodegroup_gslbvserver_binding state %s", bindingId)

	findParams := service.FindParams{
		ResourceType:             "clusternodegroup_gslbvserver_binding",
		ResourceName:             name,
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
		log.Printf("[WARN] citrixadc-provider: Clearing clusternodegroup_gslbvserver_binding state %s", bindingId)
		d.SetId("")
		return nil
	}

	// Iterate through results to find the one with the right id
	foundIndex := -1
	for i, v := range dataArr {
		if v["vserver"].(string) == vserver {
			foundIndex = i
			break
		}
	}

	// Resource is missing
	if foundIndex == -1 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams vserver not found in array")
		log.Printf("[WARN] citrixadc-provider: Clearing clusternodegroup_gslbvserver_binding state %s", bindingId)
		d.SetId("")
		return nil
	}
	// Fallthrough

	data := dataArr[foundIndex]

	d.Set("name", data["name"])
	d.Set("vserver", data["vserver"])

	return nil

}

func deleteClusternodegroup_gslbvserver_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In deleteClusternodegroup_gslbvserver_bindingFunc")
	client := meta.(*NetScalerNitroClient).client

	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	name := idSlice[0]
	vserver := idSlice[1]

	args := make([]string, 0)
	args = append(args, fmt.Sprintf("vserver:%s", vserver))

	err := client.DeleteResourceWithArgs(service.Clusternodegroup_gslbvserver_binding.Type(), name, args)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
