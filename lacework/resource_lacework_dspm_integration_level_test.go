package lacework

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression for the v2.4.0 phantom-PATCH: the DSPM integration_level attribute must be
// Optional+Computed with NO schema Default. A schema Default is applied to any config that
// omits the field, which forced an in-place update (PATCH) on every DSPM integration created
// before this attribute existed (their state has an empty integration_level) — breaking
// existing deployments on provider upgrade. Computed lets an omitted value track the backend
// with no diff; the Create/Update code still defaults the value actually sent to the API.
func TestDspmIntegrationLevel_SchemaHasNoDefault(t *testing.T) {
	for _, tc := range []struct {
		name string
		res  *schema.Resource
	}{
		{"aws", resourceLaceworkAwsDspm()},
		{"azure", resourceLaceworkAzureDspm()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.res.Schema["integration_level"]
			require.NotNil(t, s)
			assert.Nil(t, s.Default, "integration_level must NOT have a schema Default (it forces a PATCH on pre-existing integrations)")
			assert.True(t, s.Computed, "integration_level must be Computed so an omitted value tracks the backend")
			assert.True(t, s.Optional, "integration_level must remain Optional")
		})
	}
}

// TestDspmIntegrationLevel_OmittedProducesNoDiff exercises the behavior the schema guards:
// an existing integration (integration_level already empty in state) whose config omits
// integration_level must plan NO change for that attribute — i.e. Terraform won't PATCH it
// on upgrade.
func TestDspmIntegrationLevel_OmittedProducesNoDiff(t *testing.T) {
	for _, tc := range []struct {
		name   string
		res    *schema.Resource
		config map[string]interface{}
	}{
		{
			"aws", resourceLaceworkAwsDspm(),
			map[string]interface{}{
				"name":               "existing",
				"storage_bucket_arn": "arn:aws:s3:::bucket",
				"regions":            []interface{}{"us-east-1"},
				"credentials":        []interface{}{map[string]interface{}{"external_id": "x", "role_arn": "arn:aws:iam::123456789012:role/r"}},
			},
		},
		{
			"azure", resourceLaceworkAzureDspm(),
			map[string]interface{}{
				"name":                "existing",
				"storage_account_url": "https://x.blob.core.windows.net/",
				"blob_container_name": "c",
				"regions":             []interface{}{"East US"},
				"credentials":         []interface{}{map[string]interface{}{"client_id": "x", "client_secret": "y"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate an integration created before integration_level existed: empty in state.
			state := &terraform.InstanceState{
				ID:         "existing-guid",
				Attributes: map[string]string{"integration_level": ""},
			}
			cfg := terraform.NewResourceConfigRaw(tc.config)

			diff, err := tc.res.Diff(context.Background(), state, cfg, nil)
			require.NoError(t, err)

			if diff != nil {
				if d, ok := diff.Attributes["integration_level"]; ok {
					assert.Falsef(t, d.New != d.Old || d.NewComputed,
						"omitted integration_level must not plan a change (old=%q new=%q computed=%v)",
						d.Old, d.New, d.NewComputed)
				}
			}
		})
	}
}
