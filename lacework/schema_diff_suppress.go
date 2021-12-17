package lacework

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func diffCaseInsensitive(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
