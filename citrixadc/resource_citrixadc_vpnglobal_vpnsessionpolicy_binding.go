package citrixadc

import (
	"github.com/citrix/adc-nitro-go/resource/config/vpn"
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"log"
	"net/url"
)

func resourceCitrixAdcVpnglobal_vpnsessionpolicy_binding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createVpnglobal_vpnsessionpolicy_bindingFunc,
		Read:          readVpnglobal_vpnsessionpolicy_bindingFunc,
		Delete:        deleteVpnglobal_vpnsessionpolicy_bindingFunc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policyname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"builtin": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"feature": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gotopriorityexpression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"groupextraction": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secondary": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func createVpnglobal_vpnsessionpolicy_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In createVpnglobal_vpnsessionpolicy_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	policyname := d.Get("policyname").(string)
	vpnglobal_vpnsessionpolicy_binding := vpn.Vpnglobalvpnsessionpolicybinding{
		Builtin:                toStringList(d.Get("builtin").([]interface{})),
		Feature:                d.Get("feature").(string),
		Gotopriorityexpression: d.Get("gotopriorityexpression").(string),
		Groupextraction:        d.Get("groupextraction").(bool),
		Policyname:             d.Get("policyname").(string),
		Priority:               d.Get("priority").(int),
		Secondary:              d.Get("secondary").(bool),
	}

	err := client.UpdateUnnamedResource(service.Vpnglobal_vpnsessionpolicy_binding.Type(), &vpnglobal_vpnsessionpolicy_binding)
	if err != nil {
		return err
	}

	d.SetId(policyname)

	err = readVpnglobal_vpnsessionpolicy_bindingFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this vpnglobal_vpnsessionpolicy_binding but we can't read it ?? %s", policyname)
		return nil
	}
	return nil
}

func readVpnglobal_vpnsessionpolicy_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] citrixadc-provider:  In readVpnglobal_vpnsessionpolicy_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	policyname := d.Id()
	log.Printf("[DEBUG] citrixadc-provider: Reading vpnglobal_vpnsessionpolicy_binding state %s", policyname)

	findParams := service.FindParams{
		ResourceType:             "vpnglobal_vpnsessionpolicy_binding",
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
		log.Printf("[WARN] citrixadc-provider: Clearing vpnglobal_vpnsessionpolicy_binding state %s", policyname)
		d.SetId("")
		return nil
	}

	// Iterate through results to find the one with the right id
	foundIndex := -1
	for i, v := range dataArr {
		if v["policyname"].(string) == policyname {
			foundIndex = i
			break
		}
	}

	// Resource is missing
	if foundIndex == -1 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams policyname IdComponent not found in array")
		log.Printf("[WARN] citrixadc-provider: Clearing vpnglobal_vpnsessionpolicy_binding state %s", policyname)
		d.SetId("")
		return nil
	}
	// Fallthrough

	data := dataArr[foundIndex]

	d.Set("builtin", data["builtin"])
	d.Set("feature", data["feature"])
	d.Set("gotopriorityexpression", data["gotopriorityexpression"])
	d.Set("groupextraction", data["groupextraction"])
	d.Set("policyname", data["policyname"])
	d.Set("priority", data["priority"])
	d.Set("secondary", data["secondary"])

	return nil

}

func deleteVpnglobal_vpnsessionpolicy_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In deleteVpnglobal_vpnsessionpolicy_bindingFunc")
	client := meta.(*NetScalerNitroClient).client

	policyname := d.Id()
	args := make([]string, 0)
	args = append(args, fmt.Sprintf("policyname:%s", policyname))
	if val, ok := d.GetOk("secondary"); ok {
		args = append(args, fmt.Sprintf("secondary:%s", url.QueryEscape(val.(string))))
	}
	if val, ok := d.GetOk("groupextraction"); ok {
		args = append(args, fmt.Sprintf("groupextraction:%s", url.QueryEscape(val.(string))))
	}

	err := client.DeleteResourceWithArgs(service.Vpnglobal_vpnsessionpolicy_binding.Type(), "", args)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
