package citrixadc

import (
	"github.com/citrix/adc-nitro-go/resource/config/lsn"
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"log"
	"net/url"
	"strings"
)

func resourceCitrixAdcLsnclient_network6_binding() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createLsnclient_network6_bindingFunc,
		Read:          readLsnclient_network6_bindingFunc,
		Delete:        deleteLsnclient_network6_bindingFunc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"clientname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network6": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"td": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func createLsnclient_network6_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In createLsnclient_network6_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	clientname := d.Get("clientname")
	network6 := d.Get("network6")
	bindingId := fmt.Sprintf("%s,%s", clientname, network6)
	lsnclient_network6_binding := lsn.Lsnclientnetwork6binding{
		Clientname: d.Get("clientname").(string),
		Network6:   d.Get("network6").(string),
		Td:         d.Get("td").(int),
	}

	err := client.UpdateUnnamedResource("lsnclient_network6_binding", &lsnclient_network6_binding)
	if err != nil {
		return err
	}

	d.SetId(bindingId)

	err = readLsnclient_network6_bindingFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this lsnclient_network6_binding but we can't read it ?? %s", bindingId)
		return nil
	}
	return nil
}

func readLsnclient_network6_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] citrixadc-provider:  In readLsnclient_network6_bindingFunc")
	client := meta.(*NetScalerNitroClient).client
	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	clientname := idSlice[0]
	network6 := idSlice[1]

	log.Printf("[DEBUG] citrixadc-provider: Reading lsnclient_network6_binding state %s", bindingId)

	findParams := service.FindParams{
		ResourceType:             "lsnclient_network6_binding",
		ResourceName:             clientname,
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
		log.Printf("[WARN] citrixadc-provider: Clearing lsnclient_network6_binding state %s", bindingId)
		d.SetId("")
		return nil
	}

	// Iterate through results to find the one with the right id
	foundIndex := -1
	for i, v := range dataArr {
		if strings.ToLower(v["network6"].(string))  == strings.ToLower(network6) {
			foundIndex = i
			break
		}
	}

	// Resource is missing
	if foundIndex == -1 {
		log.Printf("[DEBUG] citrixadc-provider: FindResourceArrayWithParams network6 not found in array")
		log.Printf("[WARN] citrixadc-provider: Clearing lsnclient_network6_binding state %s", bindingId)
		d.SetId("")
		return nil
	}
	// Fallthrough

	data := dataArr[foundIndex]

	d.Set("clientname", data["clientname"])
	d.Set("network6", data["network6"])
	d.Set("td", data["td"])

	return nil

}

func deleteLsnclient_network6_bindingFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  citrixadc-provider: In deleteLsnclient_network6_bindingFunc")
	client := meta.(*NetScalerNitroClient).client

	bindingId := d.Id()
	idSlice := strings.SplitN(bindingId, ",", 2)

	name := idSlice[0]
	network6 := idSlice[1]

	args := make([]string, 0)
	args = append(args, fmt.Sprintf("network6:%s", url.PathEscape(network6)))
	if v, ok := d.GetOk("td"); ok {
		args = append(args, fmt.Sprintf("td:%v", v.(int)))
	}

	err := client.DeleteResourceWithArgs("lsnclient_network6_binding", name, args)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
