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

const testAccContentinspectionaction_basic = `

	resource "citrixadc_contentinspectionaction" "tf_contentinspectionaction" {
		name            = "my_ci_action"
		type            = "ICAP"
		icapprofilename = "reqmod-profile"
		servername      = "vicap"
		ifserverdown    = "DROP"
	}
`
const testAccContentinspectionaction_update = `

	resource "citrixadc_contentinspectionaction" "tf_contentinspectionaction" {
		name            = "my_ci_action"
		type            = "NOINSPECTION"
	}
`

func TestAccContentinspectionaction_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentinspectionactionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContentinspectionaction_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentinspectionactionExist("citrixadc_contentinspectionaction.tf_contentinspectionaction", nil),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "name", "my_ci_action"),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "type", "ICAP"),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "icapprofilename", "reqmod-profile"),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "servername", "vicap"),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "ifserverdown", "DROP"),
				),
			},
			resource.TestStep{
				Config: testAccContentinspectionaction_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentinspectionactionExist("citrixadc_contentinspectionaction.tf_contentinspectionaction", nil),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "name", "my_ci_action"),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionaction.tf_contentinspectionaction", "type", "NOINSPECTION"),
				),
			},
		},
	})
}

func testAccCheckContentinspectionactionExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No contentinspectionaction name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource("contentinspectionaction", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("contentinspectionaction %s not found", n)
		}

		return nil
	}
}

func testAccCheckContentinspectionactionDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_contentinspectionaction" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource("contentinspectionaction", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("contentinspectionaction %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
