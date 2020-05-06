package lacework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApiToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
        data "lacework_api_token" "test" {}
        `,
				Check: resource.TestCheckResourceAttrSet(
					"data.lacework_api_token.test", "token",
				),
			},
		},
	})
}
