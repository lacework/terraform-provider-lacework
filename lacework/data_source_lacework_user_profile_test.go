package lacework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUserProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
        data "lacework_user_profile" "test" {}
        `,
				Check: resource.TestCheckResourceAttrSet(
					"data.lacework_user_profile.test", "username",
				),
			},
		},
	})
}
