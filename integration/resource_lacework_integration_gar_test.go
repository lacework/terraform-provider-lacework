package integration

// TestIntegrationGARCreate applies integration terraform:
// => '../examples/resource_lacework_integration_gar'
//
// It uses the go-sdk to verify the created integration,
// applies an update with new integration name and destroys it
//func TestIntegrationGARCreate(t *testing.T) {
//	gcreds, err := googleLoadDefaultCredentials()
//	if assert.Nil(t, err, "this test requires you to set GOOGLE_CREDENTIALS environment variable") {
//		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
//			TerraformDir: "../examples/resource_lacework_integration_gar",
//			Vars: map[string]interface{}{
//				"client_id":              gcreds.ClientID,
//				"client_email":           gcreds.ClientEmail,
//				"private_key_id":         gcreds.PrivateKeyID,
//				"non_os_package_support": true,
//			},
//			EnvVars: map[string]string{
//				"TF_VAR_private_key": gcreds.PrivateKey, // @afiune this will avoid printing secrets in our pipeline
//				"LW_API_TOKEN":       LwApiToken,
//			},
//		})
//		defer terraform.Destroy(t, terraformOptions)
//
//		// Create new Google Artifact Registry
//		create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
//		createData := GetContainerRegisteryGar(create)
//		assert.Equal(t, "Google Artifact Registry Example", createData.Data.Name)
//		assert.Equal(t, true, createData.Data.Data.NonOSPackageEval)
//
//		// Update Google Artifact Registry
//		terraformOptions.Vars["integration_name"] = "Google Artifact Registry Updated"
//		terraformOptions.Vars["non_os_package_support"] = true
//
//		update := terraform.ApplyAndIdempotent(t, terraformOptions)
//		updateData := GetContainerRegisteryGar(update)
//		assert.Equal(t, "Google Artifact Registry Updated", updateData.Data.Name)
//		assert.Equal(t, true, updateData.Data.Data.NonOSPackageEval)
//	}
//}
