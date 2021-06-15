package lacework

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastStringSlice(t *testing.T) {
	i := make([]interface{}, 3)
	i[0] = "abc"
	i[1] = "xyz"
	i[2] = "bubulubu"

	var (
		expected = []string{"abc", "xyz", "bubulubu"}
		actual   = castStringSlice(i)
	)

	assert.ElementsMatchf(t, expected, actual,
		"%s did not match expected value: %s", actual, expected,
	)
}

func TestCastAndTransformStringSlice(t *testing.T) {
	i := make([]interface{}, 3)
	i[0] = "abc"
	i[1] = "xyz"
	i[2] = "bubulubu"

	expected := []string{"ABC", "XYZ", "BUBULUBU"}
	actual := castAndTransformStringSlice(i, strings.ToUpper)

	assert.ElementsMatchf(t, expected, actual,
		"%s did not match expected value: %s", actual, expected,
	)
}

func TestCastAttributeToStringSlice(t *testing.T) {
	var (
		expected     = []string{"foo", "bar"}
		d            = resourceLaceworkIntegrationDockerHub()
		testResource = d.TestResourceData()
	)

	testResource.Set("limit_by_tags", expected)
	actual := castAttributeToStringSlice(testResource, "limit_by_tags")

	assert.ElementsMatchf(t, expected, actual,
		"%s did not match expected value: %s", actual, expected,
	)
}
