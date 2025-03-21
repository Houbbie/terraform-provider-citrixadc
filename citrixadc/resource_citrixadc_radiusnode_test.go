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
	"net/url"
)

const testAccRadiusnode_basic = `
	resource "citrixadc_radiusnode" "tf_radiusnode" {
		nodeprefix = "10.10.10.10/32"
		radkey     = "secret"
	}
`

func TestAccRadiusnode_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRadiusnodeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRadiusnode_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRadiusnodeExist("citrixadc_radiusnode.tf_radiusnode", nil),
					resource.TestCheckResourceAttr("citrixadc_radiusnode.tf_radiusnode", "nodeprefix", "10.10.10.10/32"),
				),
			},
		},
	})
}

func testAccCheckRadiusnodeExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No radiusnode name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource("radiusnode", url.PathEscape(url.QueryEscape(rs.Primary.ID)))

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("radiusnode %s not found", n)
		}

		return nil
	}
}

func testAccCheckRadiusnodeDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_radiusnode" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource("radiusnode", url.PathEscape(url.QueryEscape(rs.Primary.ID)))
		if err == nil {
			return fmt.Errorf("radiusnode %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
