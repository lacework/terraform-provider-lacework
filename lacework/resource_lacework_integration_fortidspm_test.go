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
	assert.True(t, r.Schema["tenant_level"].Computed)
	assert.True(t, r.Schema["tenant_level"].ForceNew)
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

func azureFortiDspmDataFrom(t *testing.T, raw map[string]interface{}) api.AzureFortiDspmData {
	raw["name"] = "dspm"
	raw["tenant_id"] = "tenant-1"
	raw["regions"] = []interface{}{"westus2", "eastus"}
	return azureFortiDspmData(schema.TestResourceDataRaw(t, resourceLaceworkIntegrationAzureFortiDspm().Schema, raw))
}

func TestAzureFortiDspmData(t *testing.T) {
	cases := []struct {
		name string
		raw  map[string]interface{}
		want api.AzureFortiDspmData
	}{
		{
			"subscription given",
			map[string]interface{}{"subscription_id": "sub-1"},
			api.AzureFortiDspmData{TenantID: "tenant-1", SubscriptionID: "sub-1", Regions: []string{"westus2", "eastus"},
				IntegrationLevel: api.AzureSubscriptionIntegration},
		},
		{
			"tenant_level set keeps the subscription",
			map[string]interface{}{"subscription_id": "sub-1", "tenant_level": true},
			api.AzureFortiDspmData{TenantID: "tenant-1", SubscriptionID: "sub-1", Regions: []string{"westus2", "eastus"},
				IntegrationLevel: api.AzureTenantIntegration},
		},
		{
			"no subscription",
			map[string]interface{}{},
			api.AzureFortiDspmData{TenantID: "tenant-1", Regions: []string{"westus2", "eastus"},
				IntegrationLevel: api.AzureTenantIntegration},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, azureFortiDspmDataFrom(t, c.raw), c.name)
	}
}
