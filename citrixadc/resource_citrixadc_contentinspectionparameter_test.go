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

const testAccContentinspectionparameter_basic = `

resource "citrixadc_contentinspectionparameter" "tf_contentinspectionparameter" {
	undefaction = "RESET"
	}
`
const testAccContentinspectionparameter_update = `

resource "citrixadc_contentinspectionparameter" "tf_contentinspectionparameter" {
	undefaction = "DROP"
	}
`


func TestAccContentinspectionparameter_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContentinspectionparameter_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentinspectionparameterExist("citrixadc_contentinspectionparameter.tf_contentinspectionparameter", nil),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionparameter.tf_contentinspectionparameter", "undefaction", "RESET"),
				),
			},
			resource.TestStep{
				Config: testAccContentinspectionparameter_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentinspectionparameterExist("citrixadc_contentinspectionparameter.tf_contentinspectionparameter", nil),
					resource.TestCheckResourceAttr("citrixadc_contentinspectionparameter.tf_contentinspectionparameter", "undefaction", "DROP"),
				),
			},
		},
	})
}

func testAccCheckContentinspectionparameterExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No contentinspectionparameter name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource("contentinspectionparameter", "")

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("contentinspectionparameter %s not found", n)
		}

		return nil
	}
}