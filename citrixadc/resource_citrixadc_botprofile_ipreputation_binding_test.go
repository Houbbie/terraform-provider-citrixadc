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
	"strings"
	"testing"
)

const testAccBotprofile_ipreputation_binding_basic = `
	resource "citrixadc_botprofile" "tf_botprofile" {
		name                     = "tf_botprofile"
		errorurl                 = "http://www.citrix.com"
		trapurl                  = "/http://www.citrix.com"
		comment                  = "tf_botprofile comment"
		bot_enable_white_list    = "ON"
		bot_enable_black_list    = "ON"
		bot_enable_rate_limit    = "ON"
		devicefingerprint        = "ON"
		devicefingerprintaction  = ["LOG", "RESET"]
		bot_enable_ip_reputation = "ON"
		trap                     = "ON"
		trapaction               = ["LOG", "RESET"]
		bot_enable_tps           = "ON"
	}
	resource "citrixadc_botprofile_ipreputation_binding" "tf_binding" {
		name              = citrixadc_botprofile.tf_botprofile.name
		bot_ipreputation  = "true"
		category          = "BOTNETS"
		bot_iprep_action  = ["LOG", "REDIRECT"]
		bot_bind_comment  = "TestingIpreputation"
		bot_iprep_enabled = "ON"
		logmessage        = "MessageTesting"
	}
`

const testAccBotprofile_ipreputation_binding_basic_step2 = `
	# Keep the above bound resources without the actual binding to check proper deletion
	resource "citrixadc_botprofile" "tf_botprofile" {
		name                     = "tf_botprofile"
		errorurl                 = "http://www.citrix.com"
		trapurl                  = "/http://www.citrix.com"
		comment                  = "tf_botprofile comment"
		bot_enable_white_list    = "ON"
		bot_enable_black_list    = "ON"
		bot_enable_rate_limit    = "ON"
		devicefingerprint        = "ON"
		devicefingerprintaction  = ["LOG", "RESET"]
		bot_enable_ip_reputation = "ON"
		trap                     = "ON"
		trapaction               = ["LOG", "RESET"]
		bot_enable_tps           = "ON"
	}
`

func TestAccBotprofile_ipreputation_binding_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBotprofile_ipreputation_bindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBotprofile_ipreputation_binding_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBotprofile_ipreputation_bindingExist("citrixadc_botprofile_ipreputation_binding.tf_binding", nil),
				),
			},
			resource.TestStep{
				Config: testAccBotprofile_ipreputation_binding_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBotprofile_ipreputation_bindingNotExist("citrixadc_botprofile_ipreputation_binding.tf_binding", "tf_botprofile,BOTNETS"),
				),
			},
		},
	})
}

func testAccCheckBotprofile_ipreputation_bindingExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No botprofile_ipreputation_binding id is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		client := testAccProvider.Meta().(*NetScalerNitroClient).client

		bindingId := rs.Primary.ID

		idSlice := strings.SplitN(bindingId, ",", 2)

		name := idSlice[0]
		category := idSlice[1]

		findParams := service.FindParams{
			ResourceType:             "botprofile_ipreputation_binding",
			ResourceName:             name,
			ResourceMissingErrorCode: 258,
		}
		dataArr, err := client.FindResourceArrayWithParams(findParams)

		// Unexpected error
		if err != nil {
			return err
		}

		// Iterate through results to find the one with the matching secondIdComponent
		found := false
		for _, v := range dataArr {
			if v["category"].(string) == category {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("botprofile_ipreputation_binding %s not found", n)
		}

		return nil
	}
}

func testAccCheckBotprofile_ipreputation_bindingNotExist(n string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*NetScalerNitroClient).client

		if !strings.Contains(id, ",") {
			return fmt.Errorf("Invalid id string %v. The id string must contain a comma.", id)
		}
		idSlice := strings.SplitN(id, ",", 2)

		name := idSlice[0]
		category := idSlice[1]

		findParams := service.FindParams{
			ResourceType:             "botprofile_ipreputation_binding",
			ResourceName:             name,
			ResourceMissingErrorCode: 258,
		}
		dataArr, err := client.FindResourceArrayWithParams(findParams)

		// Unexpected error
		if err != nil {
			return err
		}

		// Iterate through results to hopefully not find the one with the matching secondIdComponent
		found := false
		for _, v := range dataArr {
			if v["category"].(string) == category {
				found = true
				break
			}
		}

		if found {
			return fmt.Errorf("botprofile_ipreputation_binding %s was found, but it should have been destroyed", n)
		}

		return nil
	}
}

func testAccCheckBotprofile_ipreputation_bindingDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_botprofile_ipreputation_binding" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource("botprofile_ipreputation_binding", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("botprofile_ipreputation_binding %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
