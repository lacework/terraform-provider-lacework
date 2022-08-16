package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lacework/go-sdk/api"
)

func dataSourceLaceworkUserProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLaceworkUserProfileRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_account": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"org_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"org_user": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accounts": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cust_guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLaceworkUserProfileRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	response, err := lacework.V2.UserProfile.Get()
	if err != nil {
		// return the api client error directly since it is user friendly
		return err
	}

	for _, profile := range response.Data {
		accounts := make([]map[string]interface{}, 0, len(profile.Accounts))
		for _, a := range profile.Accounts {
			account := make(map[string]interface{})
			account["account_name"] = a.AccountName
			account["admin"] = a.Admin
			account["cust_guid"] = a.CustGUID
			account["user_enabled"] = a.UserEnabled
			account["user_guid"] = a.UserGUID

			accounts = append(accounts, account)
		}

		d.SetId(profile.Username)
		d.Set("org_account", profile.OrgAccount)
		d.Set("org_admin", profile.OrgAdmin)
		d.Set("org_user", profile.OrgUser)
		d.Set("url", profile.URL)
		d.Set("accounts", accounts)
	}

	return nil
}
