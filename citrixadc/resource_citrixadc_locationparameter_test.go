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

const testAccLocationparameter_add = `
	resource "citrixadc_locationparameter" "tf_locationpara" {
		context            = "geographic"
		q1label            = "asia"
		matchwildcardtoany = "YES"
	}
`
const testAccLocationparameter_update = `
	resource "citrixadc_locationparameter" "tf_locationpara" {
		context            = "geographic"
		q1label            = "europe"
		matchwildcardtoany = "NO"
	}
`

func TestAccLocationparameter_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccLocationparameter_add,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationparameterExist("citrixadc_locationparameter.tf_locationpara", nil),
					resource.TestCheckResourceAttr("citrixadc_locationparameter.tf_locationpara", "context", "geographic"),
					resource.TestCheckResourceAttr("citrixadc_locationparameter.tf_locationpara", "q1label", "asia"),
					resource.TestCheckResourceAttr("citrixadc_locationparameter.tf_locationpara", "matchwildcardtoany", "YES"),
				),
			},
			resource.TestStep{
				Config: testAccLocationparameter_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationparameterExist("citrixadc_locationparameter.tf_locationpara", nil),
					resource.TestCheckResourceAttr("citrixadc_locationparameter.tf_locationpara", "context", "geographic"),
					resource.TestCheckResourceAttr("citrixadc_locationparameter.tf_locationpara", "matchwildcardtoany", "NO"),
				),
			},
		},
	})
}

func testAccCheckLocationparameterExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No locationparameter name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(service.Locationparameter.Type(), "")

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("locationparameter %s not found", n)
		}

		return nil
	}
}
