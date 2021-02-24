package lacework

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProvider  *schema.Provider
	testAccProviders map[string]*schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"lacework": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LW_API_KEY"); v == "" {
		t.Fatal("LW_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("LW_API_SECRET"); v == "" {
		t.Fatal("LW_API_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("LW_ACCOUNT"); v == "" {
		t.Fatal("LW_ACCOUNT must be set for acceptance tests")
	}
}
