package lacework

import (
	"testing"

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
