package lacework

import (
	"fmt"
	"log"
	"net/mail"
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
				Type:          schema.TypeBool,
				Default:       false,
				Optional:      true,
				ConflictsWith: []string{"organization"},
				Description:   "Set to true to make the team member an administrator, otherwise the member will be a regular user",
			},
			"organization": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"administrator": {
							Type:          schema.TypeBool,
							Optional:      true,
							Default:       false,
							ConflictsWith: []string{"organization.0.user"},
							Description:   "Whether the team member is an org level administrator",
						},
						"user": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the team member is an org level user",
						},
						"admin_accounts": {
							Type:             schema.TypeSet,
							Optional:         true,
							Description:      "List of accounts the team member is an admin",
							ConflictsWith:    []string{"organization.0.user", "organization.0.administrator"},
							DiffSuppressFunc: diffCaseInsensitive,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(val interface{}) string {
									return strings.ToUpper(strings.TrimSpace(val.(string)))
								},
							},
						},
						"user_accounts": {
							Type:             schema.TypeSet,
							Optional:         true,
							Description:      "List of accounts the team member is a user",
							ConflictsWith:    []string{"organization.0.user", "organization.0.administrator"},
							DiffSuppressFunc: diffCaseInsensitive,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(val interface{}) string {
									return strings.ToUpper(strings.TrimSpace(val.(string)))
								},
							},
						},
					},
				},
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
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

	if !d.Get("enabled").(bool) {
		tmOrg.UserEnabled = 0
	}

	tmOrg.OrgAdmin = d.Get("organization.0.administrator").(bool)
	tmOrg.OrgUser = d.Get("organization.0.user").(bool)
	tmOrg.AdminRoleAccounts = castAndUpperStringSlice(d.Get("organization.0.admin_accounts").(*schema.Set).List())
	tmOrg.UserRoleAccounts = castAndUpperStringSlice(d.Get("organization.0.user_accounts").(*schema.Set).List())

	if len(tmOrg.AdminRoleAccounts) != 0 || len(tmOrg.UserRoleAccounts) != 0 {
		// if admin_accounts or user_accounts are set, turn off OrgUser which is turned on by default
		tmOrg.OrgUser = false
	}

	if err := validateOrgTeamMember(&tmOrg); err != nil {
		return err
	}

	log.Printf("[INFO] Creating org team member with data %+v\n", tmOrg)
	response, err := lacework.V2.TeamMembers.CreateOrg(tmOrg)
	if err != nil {
		return err
	}

	if len(response.Data.Accounts) == 0 {
		msg := `
Unable to read sever response data. (empty 'accounts' field)

This was an unexpected behavior, verify that your team member has been
created successfully and report this issue to support@lacework.net
`
		return errors.New(msg)
	}

	d.SetId(response.Data.Accounts[0].UserGuid)
	d.Set("guid", response.Data.Accounts[0].UserGuid)

	org := make(map[string]interface{})
	org["admin_accounts"] = tmOrg.AdminRoleAccounts
	org["user_accounts"] = tmOrg.UserRoleAccounts
	org["user"] = tmOrg.OrgUser
	org["administrator"] = tmOrg.OrgAdmin
	log.Printf("[INFO] (Create) Setting up organization state: %v\n", org)
	d.Set("organization", []map[string]interface{}{org})

	log.Printf("[INFO] Created org team member with email %s and guid %s\n",
		response.Data.UserName, d.Id())
	return nil
}

func validateOrgTeamMember(m *api.TeamMemberOrg) error {
	if len(m.AdminRoleAccounts) != 0 || len(m.UserRoleAccounts) != 0 {
		// the user can't use the administrator argument with either admin_accounts or user_accounts
		if m.OrgAdmin {
			return errors.New("organization.0.admin_accounts and organization.0.user_accounts can't be used with the organization.0.administrator argument.")
		}

		// verify that an account doesn't exist in both lists
		for _, account := range m.AdminRoleAccounts {
			if ContainsStr(m.UserRoleAccounts, account) {
				return errors.Errorf("the same account can't be specified in both arguments, organization.0.admin_accounts and organization.0.user_accounts")
			}
		}
	}

	// org team members can't use AdminUser in props API field
	if m.Props.AccountAdmin {
		return errors.New("administrator argument can't be used when creating an organizational team member, use organization.0.administrator instead")
	}
	return nil
}

func laceworkTeamMemberCreate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	if _, ok := d.GetOk("organization.0"); ok {
		msg := `

To manage team members at the organization-level you need to define a Lacework
provider with the 'organization' argument set to 'true'.

    provider "lacework" {
      organization = true
    }

Refer to the resource documentation for more information:

    https://registry.terraform.io/providers/lacework/lacework/latest/docs/resources/team_member`
		return errors.New(msg)
	}

	tm := api.NewTeamMember(d.Get("email").(string),
		api.TeamMemberProps{
			AccountAdmin: d.Get("administrator").(bool),
			Company:      d.Get("company").(string),
			FirstName:    d.Get("first_name").(string),
			LastName:     d.Get("last_name").(string),
		})

	if !d.Get("enabled").(bool) {
		tm.UserEnabled = 0
	}

	fmt.Printf("[INFO] Creating team member with data %v\n", tm)
	response, err := lacework.V2.TeamMembers.Create(tm)
	if err != nil {
		return err
	}

	d.SetId(response.Data.UserGuid)
	d.Set("guid", response.Data.UserGuid)
	d.Set("created_time", response.Data.Props.CreatedTime)
	d.Set("updated_time", response.Data.Props.UpdatedTime)
	d.Set("updated_by", response.Data.Props.UpdatedBy)

	fmt.Printf("[INFO] Created team member with user guid %s\n", response.Data.UserGuid)
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
	lacework := meta.(*api.Client)

	var (
		response api.TeamMemberResponse
		email    = d.Get("email").(string)
	)
	if email == "" {
		// if the email is empty it means that the user imported the team member
		// resource using a guid, so we need to figure out the email before
		// processing the read read org operation
		log.Printf("[INFO] Retrieving email from org team member with guid %s\n", d.Id())
		if err := lacework.V2.TeamMembers.Get(d.Id(), &response); err != nil {
			return err
		}
		d.Set("email", response.Data.UserName)
		email = response.Data.UserName
	}

	log.Printf("[INFO] Reading org team member with email %s\n", email)
	tms, err := lacework.V2.TeamMembers.SearchUsername(email)
	if err != nil || len(tms.Data) == 0 {
		return errors.Wrapf(err, "unable to find team member with email %s", email)
	}

	org := make(map[string]interface{})
	if tms.Data[0].Props.OrgUser || tms.Data[0].Props.OrgAdmin {
		org["user"] = tms.Data[0].Props.OrgUser
		org["administrator"] = tms.Data[0].Props.OrgAdmin
	} else {

		resProfile, err := lacework.V2.UserProfile.Get()
		if err != nil || len(resProfile.Data) == 0 {
			// TODO better error
			return err
		}

		var (
			accountsInfo  = resProfile.Data[0]
			adminAccounts = []string{}
			userAccounts  = []string{}
		)

		for _, account := range tms.Data {
			accInfo, found := SearchAccountByGUID(&accountsInfo, account.CustGuid)
			if found {
				if account.Props.AccountAdmin {
					adminAccounts = append(adminAccounts, accInfo.AccountName)
				} else {
					userAccounts = append(userAccounts, accInfo.AccountName)
				}
			}
		}

		org["admin_accounts"] = adminAccounts
		org["user_accounts"] = userAccounts
	}

	log.Printf("[INFO] (Read) Setting up organization state: %v\n", org)
	d.Set("organization", []map[string]interface{}{org})

	if err := lacework.V2.TeamMembers.Get(tms.Data[0].UserGuid, &response); err != nil {
		return err
	}

	// Org team members are tough, we can't trust the user guid since it could change,
	// so for read operations we get one valid guid and set it as the Id to use it to
	// get more information about the team member
	d.SetId(response.Data.UserGuid)
	d.Set("email", response.Data.UserName)
	d.Set("first_name", response.Data.Props.FirstName)
	d.Set("last_name", response.Data.Props.LastName)
	d.Set("company", response.Data.Props.Company)
	d.Set("enabled", response.Data.UserEnabled == 1)
	d.Set("guid", response.Data.UserGuid)
	d.Set("created_time", response.Data.Props.CreatedTime)
	d.Set("updated_time", response.Data.Props.UpdatedTime)
	d.Set("updated_by", response.Data.Props.UpdatedBy)
	// @afiune we should NOT set the administrator argument for org team members

	log.Printf("[INFO] Read org team member with email %s and guid %s\n", response.Data.UserName, d.Id())
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
	d.Set("guid", response.Data.UserGuid)
	d.Set("email", response.Data.UserName)
	d.Set("first_name", response.Data.Props.FirstName)
	d.Set("last_name", response.Data.Props.LastName)
	d.Set("company", response.Data.Props.Company)
	d.Set("enabled", response.Data.UserEnabled == 1)
	d.Set("administrator", response.Data.Props.AccountAdmin)
	d.Set("created_time", response.Data.Props.CreatedTime)
	d.Set("updated_time", response.Data.Props.UpdatedTime)
	d.Set("updated_by", response.Data.Props.UpdatedBy)

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
	lacework := meta.(*api.Client)

	tmOrg := api.NewTeamMemberOrg(d.Get("email").(string),
		api.TeamMemberProps{
			AccountAdmin: d.Get("administrator").(bool),
			Company:      d.Get("company").(string),
			FirstName:    d.Get("first_name").(string),
			LastName:     d.Get("last_name").(string),
		})

	if !d.Get("enabled").(bool) {
		tmOrg.UserEnabled = 0
	}

	tmOrg.UserGuid = d.Id()

	if d.Get("organization.0.administrator").(bool) {
		// by default the go-sdk returns an organization user,
		// if 'administrator=true' we flip both flags
		tmOrg.OrgAdmin = true
		tmOrg.OrgUser = false
	}

	tmOrg.AdminRoleAccounts = castAndUpperStringSlice(d.Get("organization.0.admin_accounts").(*schema.Set).List())
	tmOrg.UserRoleAccounts = castAndUpperStringSlice(d.Get("organization.0.user_accounts").(*schema.Set).List())

	if len(tmOrg.AdminRoleAccounts) != 0 || len(tmOrg.UserRoleAccounts) != 0 {
		// if admin_accounts or user_accounts are set, turn off OrgUser which is turned on by default
		tmOrg.OrgUser = false
	}

	if err := validateOrgTeamMember(&tmOrg); err != nil {
		return err
	}

	log.Printf("[INFO] Updating org team member with data:\n%+v\n", tmOrg)
	response, err := lacework.V2.TeamMembers.UpdateOrg(tmOrg)
	if err != nil {
		return err
	}

	if len(response.Data.Accounts) == 0 {
		msg := `
Unable to read sever response data. (empty 'accounts' field)

This was an unexpected behavior, verify that your team member has been
updated successfully and report this issue to support@lacework.net
`
		return errors.New(msg)
	}

	// we should never override the Id() of the Terraform resource but,
	// our APIs does not have a way to track an org team member, so for
	// now we are allowing this until we change our APIs with:
	//
	// => https://lacework.atlassian.net/browse/RAIN-23992
	d.SetId(response.Data.Accounts[0].UserGuid)
	d.Set("guid", response.Data.Accounts[0].UserGuid)

	log.Printf("[INFO] Updated org team member with email %s and guid %s\n",
		response.Data.UserName, d.Id())
	return nil
}

func laceworkTeamMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	if _, ok := d.GetOk("organization.0"); ok {
		msg := `

To manage team members at the organization-level you need to define a Lacework
provider with the 'organization' argument set to 'true'.

    provider "lacework" {
      organization = true
    }

Refer to the resource documentation for more information:

    https://registry.terraform.io/providers/lacework/lacework/latest/docs/resources/team_member`
		return errors.New(msg)
	}

	tm := api.NewTeamMember(d.Get("email").(string),
		api.TeamMemberProps{
			AccountAdmin: d.Get("administrator").(bool),
			Company:      d.Get("company").(string),
			FirstName:    d.Get("first_name").(string),
			LastName:     d.Get("last_name").(string),
		})

	if !d.Get("enabled").(bool) {
		tm.UserEnabled = 0
	}

	tm.UserGuid = d.Id()

	log.Printf("[INFO] Updating team member with data:\n%+v\n", tm)
	response, err := lacework.V2.TeamMembers.Update(tm)
	if err != nil {
		return err
	}

	d.Set("created_time", response.Data.Props.CreatedTime)
	d.Set("updated_time", response.Data.Props.UpdatedTime)
	d.Set("updated_by", response.Data.Props.UpdatedBy)

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
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting org team member with the user guid: %s\n", d.Id())
	err := lacework.V2.TeamMembers.DeleteOrg(d.Id())
	if err != nil {
		// TODO(afiune): if we were unable to delete the org team member by ID, try by username
		//
		//   lacework.V2.TeamMembers.DeleteOrgByUsername(d.Get("email").(string))
		//
		// This needs to be added to the go-sdk/api
		// => https://lacework.atlassian.net/browse/ALLY-798
		return err
	}

	log.Printf("[INFO] Deleted org team member with user guid: %s\n", d.Id())
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
	lacework := meta.(*api.Client)

	// we have two ways to import a team member, the first one is mostly for
	// org team members where the user provides an email
	//
	//   terraform import lacework_team_member.example foo@example.com
	//
	if _, err := mail.ParseAddress(d.Id()); err == nil {
		// if the Id() is an email address, search for the team member
		log.Printf("[INFO] Importing Lacework team member with email: %s\n", d.Id())
		tms, err := lacework.V2.TeamMembers.SearchUsername(d.Id())
		if err != nil || len(tms.Data) == 0 {
			return nil, errors.Wrap(err, "unable to find team member with specified email")
		}
		d.Set("email", d.Id())
		d.SetId(tms.Data[0].UserGuid)
		return []*schema.ResourceData{d}, nil
	}

	// the second one is for standalone team members where the guid is provided
	log.Printf("[INFO] Importing Lacework team member with user guid: %s\n", d.Id())
	var response api.TeamMemberResponse
	if err := lacework.V2.TeamMembers.Get(d.Id(), &response); err != nil {
		// maybe the user is trying to import an org team member,
		// help and point them to the first import by email
		msg := `
unable to import Lacework team member with user guid.

When trying to import an organizational team member, you could use the email
instead of the id of the user:

    terraform import lacework_team_member.<name> user@example.com
`

		return nil, errors.Wrap(err, msg)
	}
	log.Printf("[INFO] Team member found with user guid: %s\n", response.Data.UserGuid)
	return []*schema.ResourceData{d}, nil
}

// TODO(afiune): move to the go-sdk/api client
// => https://lacework.atlassian.net/browse/ALLY-798
func SearchAccountByGUID(p *api.UserProfile, guid string) (*api.Account, bool) {
	for _, acc := range p.Accounts {
		if acc.CustGUID == guid {
			return &acc, true
		}
	}
	return nil, false
}
