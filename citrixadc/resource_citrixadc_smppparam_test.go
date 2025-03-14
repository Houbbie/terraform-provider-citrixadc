/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package citrixadc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const testAccSmppparam_basic = `

resource "citrixadc_smppparam" "tf_smppparam" {
	clientmode = "TRANSCEIVER"
	msgqueue   = "OFF"
	addrnpi    = 40
	addrton    = 40
  }
  
`
const testAccSmppparam_update = `

resource "citrixadc_smppparam" "tf_smppparam" {
	clientmode = "TRANSMITTERONLY"
	msgqueue   = "ON"
	addrnpi    = 50
	addrton    = 50
  }
  
`

func TestAccSmppparam_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSmppparam_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmppparamExist("citrixadc_smppparam.tf_smppparam", nil),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "clientmode", "TRANSCEIVER"),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "msgqueue", "OFF"),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "addrnpi", "40"),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "addrton", "40"),
				),
			},
			resource.TestStep{
				Config: testAccSmppparam_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmppparamExist("citrixadc_smppparam.tf_smppparam", nil),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "clientmode", "TRANSMITTERONLY"),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "msgqueue", "ON"),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "addrnpi", "50"),
					resource.TestCheckResourceAttr("citrixadc_smppparam.tf_smppparam", "addrton", "50"),
				),
			},
		},
	})
}

func testAccCheckSmppparamExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No smppparam name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource("smppparam", "")

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("smppparam %s not found", n)
		}

		return nil
	}
}