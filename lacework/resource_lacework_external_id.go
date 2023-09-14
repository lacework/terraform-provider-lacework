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

func resourceLaceworkExternalID() *schema.Resource {
	return &schema.Resource{
		Read:   schema.Noop,
		Create: resourceLaceworkExternalIDCreate,
		Delete: schema.Noop,
		Schema: map[string]*schema.Schema{
			"csp": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("The Cloud Service Provider (%s)", strings.Join(externalIDValidCsp, ", ")),
				ValidateFunc: validation.StringInSlice(externalIDValidCsp, false),
				Required:     true,
				ForceNew:     true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The account id from the CSP to be integrated",
				Required:    true,
				ForceNew:    true,
			},
			"v2": {
				Type:        schema.TypeString,
				Description: "Generated EID version 2 ('lweid:<csp>:<version>:<lw_tenant_name>:<aws_acct_id>:<random_string_size_10>')",
				Computed:    true,
			},
		},
	}
}

func resourceLaceworkExternalIDCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	url, err := lwdomain.New(lacework.URL())
	if err != nil {
		return errors.Wrap(err, "Unable to get the Lacework account")
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
			randomStringExternalID(10),
		}, ":")
	)

	d.SetId(EID)
	d.Set("v2", EID)

	return nil
}
