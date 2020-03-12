package lacework

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"lacework": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LACEWORK_API_KEY"); v == "" {
		t.Fatal("LACEWORK_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("LACEWORK_API_SECRET"); v == "" {
		t.Fatal("LACEWORK_API_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("LACEWORK_ACCOUNT"); v == "" {
		t.Fatal("LACEWORK_ACCOUNT must be set for acceptance tests")
	}
}
