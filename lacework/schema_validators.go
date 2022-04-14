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

// StringDoesNotHavePrefix returns a SchemaValidateFunc which validates that the
// provided value does not start with any of the chars.
func StringDoesNotHavePrefix(chars string) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		if strings.HasPrefix(v, chars) {
			errors = append(errors, fmt.Errorf("expected value of %s to not start with any of %q, got %v", k, chars, i))
			return warnings, errors
		}

		return warnings, errors
	})
}
