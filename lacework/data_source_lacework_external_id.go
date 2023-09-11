package lacework

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/lwdomain"
)

var externalIDValidCsp = []string{"aws", "google", "oci", "azure"}

func dataSourceLaceworkExternalID() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLaceworkExternalIDRead,
		Schema: map[string]*schema.Schema{
			"csp": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("The Cloud Service Provider (%s)", strings.Join(externalIDValidCsp, ", ")),
				ValidateFunc: validation.StringInSlice(externalIDValidCsp, false),
				Required:     true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The account id from the CSP to be integrated",
				Required:    true,
			},
			"v2": {
				Type:        schema.TypeString,
				Description: "Generated EID version 2 ('lweid:<csp>:<version>:<lw_tenant_name>:<aws_acct_id>:<random_string_size_10>')",
				Computed:    true,
			},
			"random_string": {
				Type:        schema.TypeString,
				Description: "A random generated string (size=10)",
				Computed:    true,
			},
		},
	}
}

func dataSourceLaceworkExternalIDRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	url, err := lwdomain.New(lacework.URL())
	if err != nil {
		return errors.Wrap(err, "Unable to get the Lacework account")
	}

	randomString := d.Get("random_string").(string)
	if randomString == "" {
		randomString = randString(10)
	}
	// EID V2 Format
	//
	//     lweid:<csp>:<version>:<lw_tenant_name>:<aws_acct_id>:<random_string_size_10>
	//
	var (
		version = "v2"
		EID     = strings.Join([]string{
			"lweid",
			d.Get("csp").(string),
			version,
			url.Account,
			d.Get("account_id").(string),
			d.Get("random_string").(string),
		}, ":")
	)

	d.SetId(EID)
	d.Set("v2", EID)
	d.Set("random_string", randomString)

	return nil
}
