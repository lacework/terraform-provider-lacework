package lacework

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ValidSeverity() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(value interface{}, key string) ([]string, []error) {
		switch strings.ToLower(value.(string)) {
		case "critical", "high", "medium", "low", "info":
			return nil, nil
		default:
			return nil, []error{
				fmt.Errorf(
					"%s: can only be 'Critical', 'High', 'Medium', 'Low', 'Info'", key,
				),
			}
		}
	})
}
