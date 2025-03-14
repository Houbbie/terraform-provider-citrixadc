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
	"github.com/citrix/adc-nitro-go/service"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const testAccVpntrafficaction_add = `

	resource "citrixadc_vpntrafficaction" "tf_action" {
		name       = "Testing"
		qual       = "tcp"
		apptimeout = 20
		fta        = "OFF"
		hdx        = "OFF"
		sso        = "ON"
	}
`
const testAccVpntrafficaction_update = `

	resource "citrixadc_vpntrafficaction" "tf_action" {
		name       = "Testing"
		qual       = "tcp"
		apptimeout = 30
		fta        = "OFF"
		hdx        = "OFF"
		sso        = "OFF"
	}
`

func TestAccVpntrafficaction_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpntrafficactionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpntrafficaction_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpntrafficactionExist("citrixadc_vpntrafficaction.tf_action", nil),
					resource.TestCheckResourceAttr("citrixadc_vpntrafficaction.tf_action", "name", "Testing"),
					resource.TestCheckResourceAttr("citrixadc_vpntrafficaction.tf_action", "apptimeout", "20"),
					resource.TestCheckResourceAttr("citrixadc_vpntrafficaction.tf_action", "sso", "ON"),
				),
			},
			resource.TestStep{
				Config: testAccVpntrafficaction_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpntrafficactionExist("citrixadc_vpntrafficaction.tf_action", nil),
					resource.TestCheckResourceAttr("citrixadc_vpntrafficaction.tf_action", "name", "Testing"),
					resource.TestCheckResourceAttr("citrixadc_vpntrafficaction.tf_action", "apptimeout", "30"),
					resource.TestCheckResourceAttr("citrixadc_vpntrafficaction.tf_action", "sso", "OFF"),
				),
			},
		},
	})
}

func testAccCheckVpntrafficactionExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No vpntrafficaction name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(service.Vpntrafficaction.Type(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("vpntrafficaction %s not found", n)
		}

		return nil
	}
}

func testAccCheckVpntrafficactionDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_vpntrafficaction" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(service.Vpntrafficaction.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("vpntrafficaction %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
