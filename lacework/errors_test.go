package lacework

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Tests when a 404 http response is returned, resource id is set to empty string
func TestErrorsResourceNotFound(t *testing.T) {
	var d schema.ResourceData
	d.SetId("MockID")
	assert.Equal(t, d.Id(), "MockID")
	notFound := errors.New(`[GET] https://customerdemo.lacework.net/api/v2/VulnerabilityExceptions/VULN_ABCD" [404] Not found`)
	err := resourceNotFound(&d, notFound)
	assert.NoError(t, err)
	assert.Equal(t, d.Id(), "")
}
