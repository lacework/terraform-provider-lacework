package lacework

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
)

func resourceLaceworkTeamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkTeamMemberCreate,
		Read:   resourceLaceworkTeamMemberRead,
		Update: resourceLaceworkTeamMemberUpdate,
		Delete: resourceLaceworkTeamMemberDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkTeamMember,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email for the team member which will also be used as the username",
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The first name of the team member",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The last name of the team member",
			},
			"company": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The company name",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "The state of the team member, whether they are enabled or not",
			},
			"administrator": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the team member has admin role access into the Lacework account",
			},
			"organization": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"administrator": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the team member is an admin at the org level for the account",
						},
						"user": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the team member is an org level user",
						},
						"admin_accounts": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(val interface{}) string {
									return strings.TrimSpace(val.(string))
								},
							},
							Optional:    true,
							Description: "List of accounts the team member is an admin in",
						},
						"user_accounts": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(val interface{}) string {
									return strings.TrimSpace(val.(string))
								},
							},
							Optional:    true,
							Description: "List of accounts the team member is a user in",
						},
					},
				},
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkTeamMemberCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)
	if lacework.OrgAccess() {
		return laceworkTeamMemberCreateOrg(d, meta)
	}
	return laceworkTeamMemberCreate(d, meta)
}

func laceworkTeamMemberCreateOrg(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	tmOrg := api.NewTeamMemberOrg(d.Get("email").(string),
		api.TeamMemberProps{
			AccountAdmin: d.Get("administrator").(bool),
			Company:      d.Get("company").(string),
			FirstName:    d.Get("first_name").(string),
			LastName:     d.Get("last_name").(string),
		})

	var enabled int
	if d.Get("enabled").(bool) {
		enabled = 1
	}
	tmOrg.UserEnabled = enabled

	// Validate that the user isn't trying to be both an admin and org user
	orgAdmin := d.Get("organization.0.administrator").(bool)
	orgUser := d.Get("organization.0.user").(bool)
	if orgAdmin && orgUser {
		return errors.New("team member cannot be both an admin or an org and a user of an org")
	}
	tmOrg.OrgAdmin = orgAdmin
	tmOrg.OrgUser = orgUser
	adminAccounts := castStringSlice(d.Get("organization.0.admin_accounts").([]interface{}))
	userAccounts := castStringSlice(d.Get("organization.0.user_accounts").([]interface{}))

	var upperAdminAccounts []string
	if len(adminAccounts) > 0 {
		for _, adminAccount := range adminAccounts {
			upperAdminAccounts = append(upperAdminAccounts, strings.ToUpper(adminAccount))
		}
	}

	var upperUserAccounts []string
	if len(userAccounts) > 0 {
		for _, userAccount := range userAccounts {
			upperUserAccounts = append(upperUserAccounts, strings.ToUpper(userAccount))
		}
	}

	tmOrg.AdminRoleAccounts = upperAdminAccounts
	tmOrg.UserRoleAccounts = userAccounts

	log.Printf("[Info] Creating org team member with data %v\n", tmOrg)
	response, err := lacework.V2.TeamMembers.CreateOrg(tmOrg)
	if err != nil {
		return err
	}

	d.SetId(response.Data.UserName)
	d.Set("email", response.Data.UserName)
	d.Set("guid", response.Data.UserName)

	log.Printf("[INF0] Created org team member with username %s\n", response.Data.UserName)
	return nil
}

func laceworkTeamMemberCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	tm := api.NewTeamMember(d.Get("email").(string),
		api.TeamMemberProps{
			AccountAdmin: d.Get("administrator").(bool),
			Company:      d.Get("company").(string),
			FirstName:    d.Get("first_name").(string),
			LastName:     d.Get("last_name").(string),
		})

	var enabled int
	if d.Get("enabled").(bool) {
		enabled = 1
	}
	tm.UserEnabled = enabled

	fmt.Printf("[INFO] Creating team member with data %v\n", tm)
	response, err := lacework.V2.TeamMembers.Create(tm)
	if err != nil {
		return err
	}

	d.SetId(response.Data.UserGuid)
	d.Set("email", response.Data.UserName)
	d.Set("first_name", response.Data.Props.FirstName)
	d.Set("last_name", response.Data.Props.LastName)
	d.Set("company", response.Data.Props.Company)
	d.Set("enabled", response.Data.UserEnabled == 1)
	d.Set("administrator", response.Data.Props.AccountAdmin)
	d.Set("guid", response.Data.UserGuid)

	fmt.Printf("[INF0] Created team member with user guid %s\n", response.Data.UserGuid)
	return nil
}

func resourceLaceworkTeamMemberRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	if lacework.OrgAccess() {
		return laceworkTeamMemberReadOrg(d, meta)
	}

	return laceworkTeamMemberRead(d, meta)
}

func laceworkTeamMemberReadOrg(d *schema.ResourceData, meta interface{}) error {
	// TODO implement me please and thank you :)
	return nil
}

func laceworkTeamMemberRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading team member with user guid %s\n", d.Id())

	var response api.TeamMemberResponse
	if err := lacework.V2.TeamMembers.Get(d.Id(), &response); err != nil {
		return err
	}

	d.SetId(response.Data.UserGuid)
	d.Set("email", response.Data.UserName)
	d.Set("first_name", response.Data.Props.FirstName)
	d.Set("last_name", response.Data.Props.LastName)
	d.Set("company", response.Data.Props.Company)
	d.Set("enabled", response.Data.UserEnabled == 1)
	d.Set("administrator", response.Data.Props.AccountAdmin)
	d.Set("guid", response.Data.UserGuid)
	// The organization information here should be empty because this is an account level read
	org := make(map[string]interface{})
	d.Set("organization", []map[string]interface{}{org})

	log.Printf("[INFO] Read team member with user guid %s\n", response.Data.UserGuid)
	return nil
}

func resourceLaceworkTeamMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	if lacework.OrgAccess() {
		return laceworkTeamMemberUpdateOrg(d, meta)
	}

	return laceworkTeamMemberUpdate(d, meta)
}

func laceworkTeamMemberUpdateOrg(d *schema.ResourceData, meta interface{}) error {
	// TODO implement me please and thank you :)
	return nil
}

func laceworkTeamMemberUpdate(d *schema.ResourceData, meta interface{}) error {

	lacework := meta.(*api.Client)

	tm := api.NewTeamMember(d.Get("email").(string),
		api.TeamMemberProps{
			AccountAdmin: d.Get("administrator").(bool),
			Company:      d.Get("company").(string),
			FirstName:    d.Get("first_name").(string),
			LastName:     d.Get("last_name").(string),
		})

	var enabled int
	if d.Get("enabled").(bool) {
		enabled = 1
	}
	tm.UserEnabled = enabled
	tm.UserGuid = d.Id()

	log.Printf("[INFO] Updating team member with data:\n%+v\n", tm)
	response, err := lacework.V2.TeamMembers.Update(tm)
	if err != nil {
		return err
	}

	d.SetId(response.Data.UserGuid)
	d.Set("email", response.Data.UserName)
	d.Set("first_name", response.Data.Props.FirstName)
	d.Set("last_name", response.Data.Props.LastName)
	d.Set("company", response.Data.Props.Company)
	d.Set("enabled", response.Data.UserEnabled == 1)
	d.Set("administrator", response.Data.Props.AccountAdmin)
	d.Set("guid", response.Data.UserGuid)
	// The organization information here should be empty because this is an account level read
	org := make(map[string]interface{})
	d.Set("organization", []map[string]interface{}{org})

	log.Printf("[INFO] Updated team member with user guid %s\n", response.Data.UserGuid)
	return nil

}

func resourceLaceworkTeamMemberDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	if lacework.OrgAccess() {
		return laceworkTeamMemberDeleteOrg(d, meta)
	}
	return laceworkTeamMemberDelete(d, meta)
}

func laceworkTeamMemberDeleteOrg(d *schema.ResourceData, meta interface{}) error {
	// TODO implement me please and thank you :)
	return nil
}

func laceworkTeamMemberDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting team member with the user guid: %s\n", d.Id())
	err := lacework.V2.TeamMembers.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted team member with user guid: %s\n", d.Id())
	return nil
}

func importLaceworkTeamMember(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var response api.TeamMemberResponse
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Importing Lacework Team Member with user guid: %s\n", d.Id())

	if err := lacework.V2.TeamMembers.Get(d.Id(), &response); err != nil {
		return nil,
			errors.Wrapf(err, "unable to import Lacework resource. Team member with user guid '%s' was not found", d.Id())
	}
	log.Printf("[INFO] Team Member with user guid: %s\n found", response.Data.UserGuid)
	return []*schema.ResourceData{d}, nil
}

func searchAccountNames(d *schema.ResourceData, meta interface{}, username string) (adminAccountNames, userAccountNames []string, err error) {
	lacework := meta.(*api.Client)

	// Search for the user by username to get a list of all the account names
	tmOrgSearch, err := lacework.V2.TeamMembers.SearchUsername(username)
	if err != nil {
		return
	}

	var orgAccountGuids []string
	var userAccountGuids []string

	for _, tmOrgAccount := range tmOrgSearch.Data {
		if tmOrgAccount.Props.OrgUser {
			userAccountGuids = append(userAccountGuids, tmOrgAccount.CustGuid)
		}
		if tmOrgAccount.Props.OrgAdmin {
			orgAccountGuids = append(orgAccountGuids, tmOrgAccount.CustGuid)
		}
	}

	for _, adminAccountGuid := range orgAccountGuids {
		res, err := lacework.V2.UserProfile.Get()
		if err != nil {
			return
		}
		for _, userProfile := range res.Data {
			for _, account := range userProfile.Accounts {
				if account.CustGUID == adminAccountGuid {
					adminAccountNames = append(adminAccountNames, strings.ToUpper(account.AccountName))
				}
			}
		}
	}

	for _, userAccountGuid := range userAccountGuids {
		res, profileErr := lacework.V2.UserProfile.Get()
		if profileErr != nil {
			err = profileErr
			return
		}
		for _, userProfile := range res.Data {
			for _, account := range userProfile.Accounts {
				if account.CustGUID == userAccountGuid {
					userAccountNames = append(userAccountNames, strings.ToUpper(account.AccountName))
				}
			}
		}
	}
	return

}