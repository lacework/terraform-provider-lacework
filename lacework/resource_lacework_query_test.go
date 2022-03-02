package lacework

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIDFromLQLQuery(t *testing.T) {
	cases := []struct {
		Query           string
		expectedQueryID string
		expectedError   string
	}{
		{Query: "MyLQLQuery                   {}",
			expectedQueryID: "MyLQLQuery"},
		{Query: "foo123{source {bar} return distinct {xyz}}",
			expectedQueryID: "foo123"},
		{Query: `

     ASuper_Long_ID_madeUP_with_Numbers_1234

{
source {bar}
filter {abc}
return distinct {xyz}
}`,
			expectedQueryID: "ASuper_Long_ID_madeUP_with_Numbers_1234"},
		{Query: `Lql_Terraform_Query{
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
        and ERROR_CODE is null
    }
    return distinct {
        INSERT_ID,
        INSERT_TIME,
        EVENT_TIME,
        EVENT
    }
}`,
			expectedQueryID: "Lql_Terraform_Query"},

		// Errors!!!!
		{Query: "",
			expectedError: `query id not found. (malformed)

> Your query:


> Compare provided query to the example at:

    https://docs.lacework.com/lql-overview
`},
		{Query: "{}",
			expectedError: `query id not found. (malformed)

> Your query:
{}

> Compare provided query to the example at:

    https://docs.lacework.com/lql-overview
`},
		{Query: "     {}",
			expectedError: `query id not found. (malformed)

> Your query:
     {}

> Compare provided query to the example at:

    https://docs.lacework.com/lql-overview
`},
		{Query: `MyQueryWrongFormat }`,
			expectedError: `query id not found. (malformed)

> Your query:
MyQueryWrongFormat }

> Compare provided query to the example at:

    https://docs.lacework.com/lql-overview
`},
		{Query: `MyQueryWrongFormat } 
{
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
        and ERROR_CODE is null
    }
    return distinct {
        INSERT_ID,
        INSERT_TIME,
        EVENT_TIME,
        EVENT
    }
}`,
			expectedError: `query id not found. (malformed)

> Your query:
MyQueryWrongFormat } 
{
    source {
        CloudTrailRawEvents
    }
    filter {
        EVENT_SOURCE = 'signin.amazonaws.com'
        and EVENT_NAME in ('ConsoleLogin')
        and EVENT:additionalEventData.MFAUsed::String = 'No'
        and EVENT:responseElements.ConsoleLogin::String = 'Success'
        and ERROR_CODE is null
    }
    return distinct {
        INSERT_ID,
        INSERT_TIME,
        EVENT_TIME,
        EVENT
    }
}

> Compare provided query to the example at:

    https://docs.lacework.com/lql-overview
`},
	}

	for i, kase := range cases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actualID, err := getIDFromLQLQuery(kase.Query)
			if kase.expectedError != "" {
				if assert.Error(t, err, "should have failed") {
					assert.Equal(t, kase.expectedError, err.Error())
				}
			}
			assert.Equal(t, kase.expectedQueryID, actualID)
		})
	}
}
