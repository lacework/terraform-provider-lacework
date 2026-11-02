package lacework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
	"github.com/stretchr/testify/assert"
)

func TestResourceLaceworkIntegrationAwsFortiDspmSchema(t *testing.T) {
	r := resourceLaceworkIntegrationAwsFortiDspm()
	assert.Nil(t, r.InternalValidate(nil, true))
	assert.True(t, r.Schema["activation_tokens"].Sensitive)
	assert.True(t, r.Schema["regions"].ForceNew)
	assert.True(t, r.Schema["account_id"].ForceNew)
	for _, computed := range []string{"intg_guid", "deployment_id", "deployment_name", "env_id",
		"token_expires_in", "activation_tokens", "image_ids"} {
		assert.True(t, r.Schema[computed].Computed, computed)
	}
}

func TestResourceLaceworkIntegrationAzureFortiDspmSchema(t *testing.T) {
	r := resourceLaceworkIntegrationAzureFortiDspm()
	assert.Nil(t, r.InternalValidate(nil, true))
	assert.True(t, r.Schema["activation_tokens"].Sensitive)
	assert.True(t, r.Schema["image_urls"].Sensitive)
	assert.True(t, r.Schema["tenant_id"].ForceNew)
	assert.False(t, r.Schema["subscription_id"].Required)
	for _, computed := range []string{"intg_guid", "deployment_id", "image_url_expires_in",
		"activation_tokens", "image_urls", "hyperv_generations"} {
		assert.True(t, r.Schema[computed].Computed, computed)
	}
}

func TestResourceLaceworkFortiDspmDeploymentStatusSchema(t *testing.T) {
	r := resourceLaceworkFortiDspmDeploymentStatus()
	assert.Nil(t, r.InternalValidate(nil, true))
	assert.True(t, r.Schema["intg_guid"].ForceNew)
	assert.Equal(t, "succeeded", r.Schema["status"].Default)
	assert.Equal(t, 1, r.Schema["region"].MinItems)
}

func TestAzureFortiDspmIntegrationLevel(t *testing.T) {
	cases := []struct {
		name           string
		tenantLevel    bool
		subscriptionID string
		want           string
	}{
		{"subscription given", false, "sub-1", api.AzureSubscriptionIntegration},
		{"tenant_level set with a subscription", true, "sub-1", api.AzureTenantIntegration},
		{"tenant_level set without a subscription", true, "", api.AzureTenantIntegration},
		{"no subscription", false, "", api.AzureTenantIntegration},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, azureFortiDspmIntegrationLevel(c.tenantLevel, c.subscriptionID), c.name)
	}
}

func TestAzureFortiDspmData(t *testing.T) {
	r := resourceLaceworkIntegrationAzureFortiDspm()

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":            "dspm",
		"tenant_id":       "tenant-1",
		"subscription_id": "sub-1",
		"regions":         []interface{}{"westus2", "eastus"},
	})
	assert.Equal(t, api.AzureFortiDspmData{
		TenantID:         "tenant-1",
		SubscriptionID:   "sub-1",
		Regions:          []string{"westus2", "eastus"},
		IntegrationLevel: api.AzureSubscriptionIntegration,
	}, azureFortiDspmData(d))

	d = schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":            "dspm",
		"tenant_id":       "tenant-1",
		"subscription_id": "sub-1",
		"tenant_level":    true,
		"regions":         []interface{}{"westus2"},
	})
	assert.Equal(t, api.AzureTenantIntegration, azureFortiDspmData(d).IntegrationLevel)

	d = schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":      "dspm",
		"tenant_id": "tenant-1",
		"regions":   []interface{}{"westus2"},
	})
	data := azureFortiDspmData(d)
	assert.Equal(t, "", data.SubscriptionID)
	assert.Equal(t, api.AzureTenantIntegration, data.IntegrationLevel)
}
