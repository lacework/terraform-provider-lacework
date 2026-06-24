package lacework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// awsDspmData / azureDspmData build a *schema.ResourceData from the real resource
// schema with no backend — the cheap unit layer for the DSPM props helpers.
func awsDspmData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, resourceLaceworkAwsDspm().Schema, raw)
}

func azureDspmData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, resourceLaceworkAzureDspm().Schema, raw)
}

// ---------------------------------------------------------------------------
// buildDspmProps
// ---------------------------------------------------------------------------

func TestBuildDspmProps_NilWhenNothingSet(t *testing.T) {
	props, err := buildDspmProps(awsDspmData(t, map[string]interface{}{}), "account_filters", "account_ids")
	require.NoError(t, err)
	assert.Nil(t, props, "no optional fields set should yield nil props")
}

func TestBuildDspmProps_ScalarFields(t *testing.T) {
	d := awsDspmData(t, map[string]interface{}{
		"scan_frequency_hours": 12,
		"max_file_size_mb":     30,
	})
	props, err := buildDspmProps(d, "account_filters", "account_ids")
	require.NoError(t, err)
	require.NotNil(t, props)
	require.NotNil(t, props.Dspm.ScanIntervalHours)
	assert.Equal(t, 12, *props.Dspm.ScanIntervalHours)
	require.NotNil(t, props.Dspm.MaxDownloadBytes)
	assert.Equal(t, 30*1024*1024, *props.Dspm.MaxDownloadBytes)
}

func TestBuildDspmProps_DatastoreFilters(t *testing.T) {
	t.Run("INCLUDE with names", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"datastore_filters": []interface{}{map[string]interface{}{
				"filter_mode":     "INCLUDE",
				"datastore_names": []interface{}{"bucket-a", "bucket-b"},
			}},
		})
		props, err := buildDspmProps(d, "account_filters", "account_ids")
		require.NoError(t, err)
		require.NotNil(t, props.Dspm.DatastoreFilters)
		assert.Equal(t, "INCLUDE", props.Dspm.DatastoreFilters.FilterMode)
		assert.Equal(t, []string{"bucket-a", "bucket-b"}, props.Dspm.DatastoreFilters.DatastoreNames)
	})

	t.Run("ALL with no names is allowed", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"datastore_filters": []interface{}{map[string]interface{}{"filter_mode": "ALL"}},
		})
		props, err := buildDspmProps(d, "account_filters", "account_ids")
		require.NoError(t, err)
		require.NotNil(t, props.Dspm.DatastoreFilters)
		assert.Equal(t, "ALL", props.Dspm.DatastoreFilters.FilterMode)
		assert.Empty(t, props.Dspm.DatastoreFilters.DatastoreNames)
	})

	t.Run("ALL with names is rejected", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"datastore_filters": []interface{}{map[string]interface{}{
				"filter_mode":     "ALL",
				"datastore_names": []interface{}{"bucket-a"},
			}},
		})
		_, err := buildDspmProps(d, "account_filters", "account_ids")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be set when filter_mode is 'ALL'")
	})

	t.Run("INCLUDE with no names is rejected", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"datastore_filters": []interface{}{map[string]interface{}{"filter_mode": "INCLUDE"}},
		})
		_, err := buildDspmProps(d, "account_filters", "account_ids")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "datastore_names is required")
	})
}

func TestBuildDspmProps_AwsAccountFilters(t *testing.T) {
	t.Run("INCLUDE with account ids", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"account_filters": []interface{}{map[string]interface{}{
				"filter_mode": "INCLUDE",
				"account_ids": []interface{}{"111111111111", "222222222222"},
			}},
		})
		props, err := buildDspmProps(d, "account_filters", "account_ids")
		require.NoError(t, err)
		require.NotNil(t, props.Dspm.AccountFilters)
		assert.Equal(t, "INCLUDE", props.Dspm.AccountFilters.FilterMode)
		assert.Equal(t, []string{"111111111111", "222222222222"}, props.Dspm.AccountFilters.AccountIds)
	})

	// ALL means "scan everything", so no account ids are required (PR review fix).
	t.Run("ALL with no account ids is allowed", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"account_filters": []interface{}{map[string]interface{}{"filter_mode": "ALL"}},
		})
		props, err := buildDspmProps(d, "account_filters", "account_ids")
		require.NoError(t, err)
		require.NotNil(t, props.Dspm.AccountFilters)
		assert.Equal(t, "ALL", props.Dspm.AccountFilters.FilterMode)
		assert.Empty(t, props.Dspm.AccountFilters.AccountIds)
	})

	t.Run("INCLUDE with no account ids is rejected", func(t *testing.T) {
		d := awsDspmData(t, map[string]interface{}{
			"account_filters": []interface{}{map[string]interface{}{"filter_mode": "INCLUDE"}},
		})
		_, err := buildDspmProps(d, "account_filters", "account_ids")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "account_ids is required")
	})
}

// The shared helper is parameterized by HCL key so Azure's subscription_filters /
// subscription_ids map to the same cloud-agnostic ACCOUNT_FILTERS prop.
func TestBuildDspmProps_AzureSubscriptionFiltersMapToAccountFilters(t *testing.T) {
	d := azureDspmData(t, map[string]interface{}{
		"subscription_filters": []interface{}{map[string]interface{}{
			"filter_mode":      "EXCLUDE",
			"subscription_ids": []interface{}{"sub-a", "sub-b"},
		}},
	})
	props, err := buildDspmProps(d, "subscription_filters", "subscription_ids")
	require.NoError(t, err)
	require.NotNil(t, props.Dspm.AccountFilters)
	assert.Equal(t, "EXCLUDE", props.Dspm.AccountFilters.FilterMode)
	assert.Equal(t, []string{"sub-a", "sub-b"}, props.Dspm.AccountFilters.AccountIds)
}

// ---------------------------------------------------------------------------
// readDspmProps
// ---------------------------------------------------------------------------

func TestReadDspmProps_AwsAccountFilters(t *testing.T) {
	d := awsDspmData(t, map[string]interface{}{})
	hours, bytes := 12, 30*1024*1024
	readDspmProps(d, &api.DspmProps{Dspm: &api.DspmPropsConfig{
		ScanIntervalHours: &hours,
		MaxDownloadBytes:  &bytes,
		AccountFilters:    &api.DspmAccountFilters{FilterMode: "INCLUDE", AccountIds: []string{"111111111111"}},
		DatastoreFilters:  &api.DspmDatastoreFilters{FilterMode: "ALL"},
	}}, "account_filters", "account_ids")

	assert.Equal(t, 12, d.Get("scan_frequency_hours"))
	assert.Equal(t, 30, d.Get("max_file_size_mb"))

	af := d.Get("account_filters").([]interface{})
	require.Len(t, af, 1)
	afMap := af[0].(map[string]interface{})
	assert.Equal(t, "INCLUDE", afMap["filter_mode"])
	assert.Equal(t, []interface{}{"111111111111"}, afMap["account_ids"])
}

// readDspmProps writes the account filter under the resource's own HCL key:
// Azure should populate subscription_filters / subscription_ids.
func TestReadDspmProps_AzureSubscriptionFilters(t *testing.T) {
	d := azureDspmData(t, map[string]interface{}{})
	readDspmProps(d, &api.DspmProps{Dspm: &api.DspmPropsConfig{
		AccountFilters: &api.DspmAccountFilters{FilterMode: "EXCLUDE", AccountIds: []string{"sub-a"}},
	}}, "subscription_filters", "subscription_ids")

	sf := d.Get("subscription_filters").([]interface{})
	require.Len(t, sf, 1)
	sfMap := sf[0].(map[string]interface{})
	assert.Equal(t, "EXCLUDE", sfMap["filter_mode"])
	assert.Equal(t, []interface{}{"sub-a"}, sfMap["subscription_ids"])
}

func TestReadDspmProps_NilIsNoOp(t *testing.T) {
	d := awsDspmData(t, map[string]interface{}{})
	require.NotPanics(t, func() {
		readDspmProps(d, nil, "account_filters", "account_ids")
		readDspmProps(d, &api.DspmProps{}, "account_filters", "account_ids")
	})
}

// ---------------------------------------------------------------------------
// integration_level ValidateFunc
// ---------------------------------------------------------------------------

func TestDspmIntegrationLevelValidate(t *testing.T) {
	cases := []struct {
		name   string
		schema *schema.Resource
		valid  []string
	}{
		{"aws", resourceLaceworkAwsDspm(), []string{"ACCOUNT", "ORG", "account", "org"}},
		{"azure", resourceLaceworkAzureDspm(), []string{"TENANT", "SUBSCRIPTION", "tenant", "subscription"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vf := tc.schema.Schema["integration_level"].ValidateFunc
			require.NotNil(t, vf, "integration_level should have a ValidateFunc")

			for _, v := range tc.valid {
				_, errs := vf(v, "integration_level")
				assert.Emptyf(t, errs, "%q should be valid for %s", v, tc.name)
			}

			_, errs := vf("BOGUS", "integration_level")
			assert.NotEmpty(t, errs, "an unknown level should be rejected")
		})
	}
}
