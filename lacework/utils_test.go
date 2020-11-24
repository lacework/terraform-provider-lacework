package lacework

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastStringArray(t *testing.T) {
	i := make([]interface{}, 3)
	i[0] = "abc"
	i[1] = "xyz"
	i[2] = "bubulubu"

	var (
		expected = []string{"abc", "xyz", "bubulubu"}
		actual   = castStringArray(i)
	)

	assert.ElementsMatchf(t, expected, actual,
		"%s did not match expected value: %s", actual, expected,
	)
}

func TestCastAndTransformStringArray(t *testing.T) {
	i := make([]interface{}, 3)
	i[0] = "abc"
	i[1] = "xyz"
	i[2] = "bubulubu"

	expected := []string{"ABC", "XYZ", "BUBULUBU"}
	actual := castAndTransformStringArray(i, strings.ToUpper)

	assert.ElementsMatchf(t, expected, actual,
		"%s did not match expected value: %s", actual, expected,
	)
}
