package lacework

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
	"github.com/pkg/errors"
	"log"
	"strings"
)

func resourceLaceworkTeamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkTeamMemberCreate,
		Read: resourceLaceworkTeamMemberRead,
		Update: resourceLaceworkTeamMemberUpdate,
		Delete: resourceLaceworkTeamMemberDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkTeamMember,
				},
		Schema: map[string]*schema.Schema{
			"email": {
				Type: schema.TypeString,
				Required: true,
				Description: "The email for the team member which will also be used as the username",
			},
			"first_name": {
				Type: schema.TypeString,
				Required: true,
				Description: "The first name of the team member",
			},
			"last_name": {
				Type: schema.TypeString,
				Required: true,
				Description: "The last name of the team member",
			},
			"company": {
				Type: schema.TypeString,
				Required: true,
				Description: "The company name",
			},
			"enabled": {
				Type: schema.TypeBool,
				Required: true,
				Description: "The state of the team member, whether they are enabled or not",
			},
			"administrator": {
				Type: schema.TypeBool,
				Optional: true,
				Description: "Whether the team member has admin role access into the Lacework account",
			},
			"organization": {
				Type: schema.TypeList,
				Optional:true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"administrator": {
							Type: schema.TypeBool,
							Optional: true,
							Default: false,
							Description: "Whether the team member is an admin at the org level for the account",
						},
						"user": {
							Type: schema.TypeBool,
							Optional: true,
							Default: false,
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
							Optional: true,
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
							Optional: true,
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

	var orgTeamMember bool
	var orgAdmin bool
	var orgUser bool
	var adminAccounts []string
	var userAccounts []string
	var tmOrg api.TeamMemberOrg

	if _, ok := d.GetOk("organization"); ok {
		orgAdmin = d.Get("organization.0.administrator").(bool)
		orgUser = d.Get("organization.0.user").(bool)
		adminAccounts = d.Get("organization.0.admin_accounts").([]string)
		userAccounts = d.Get("organization.0.user_accounts").([]string)
	}

	if orgAdmin || orgUser || len(adminAccounts) > 0 || len(userAccounts) > 0{
		orgTeamMember = true
	}

	username := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	company := d.Get("company").(string)
	enabled := d.Get("enabled").(bool)
	accountAdmin := d.Get("administrator").(bool)

	var userEnabled int
	if enabled {
		userEnabled = 1
	}

	if orgUser {
		tmOrg = api.TeamMemberOrg{
			UserEnabled: userEnabled,
			UserName: username,
			Props: api.TeamMemberProps{
				FirstName: firstName,
				LastName: lastName,
				Company: company,
				AccountAdmin: accountAdmin,
			},
			AdminRoleAccounts: adminAccounts,
		}
	}

	if orgAdmin {
		tmOrg = api.TeamMemberOrg{
			Props: api.TeamMemberProps{
				FirstName: firstName,
				LastName:  lastName,
				Company:   company,
				AccountAdmin: accountAdmin,
			},
			UserEnabled:      userEnabled,
			UserName:         username,
		}
	}

	if len(userAccounts) > 0 || len(adminAccounts) > 0{
		tmOrg = api.TeamMemberOrg{
			AdminRoleAccounts: adminAccounts,
			OrgAdmin:          false,
			OrgUser:           false,
			Props:             api.TeamMemberProps{
				FirstName: firstName,
				LastName: lastName,
				Company: company,
				AccountAdmin: accountAdmin,
			},
			UserEnabled:       userEnabled,
			UserName:          username,
			UserRoleAccounts:  userAccounts,
		}
	}

	tmOrg.UserGuid = d.Id()

	if orgTeamMember {
		log.Printf("[INFO] Creating org team member with ")
		response, err := lacework.V2.TeamMembers.CreateOrg(tmOrg)
		if err != nil {
			return err
		}
		d.SetId(response.Data.UserName)
		d.Set("email", response.Data.UserName)
		d.Set("guid", response.Data.Accounts[0].UserGuid)
		log.Printf("[INFO] Created team member with username %s\n", response.Data.UserName)
		return nil
	}

	var tm api.TeamMember

	tm = api.NewTeamMember(username,
		api.TeamMemberProps{
			FirstName: firstName,
			LastName: lastName,
			Company: company,
			AccountAdmin: accountAdmin,
		})

	response, err := lacework.V2.TeamMembers.Create(tm)
	if err != nil {
		return err
	}

	d.SetId(response.Data.UserGuid)
	d.Set("email", response.Data.UserName)
	d.Set("enabled", response.Data.UserEnabled == 1)
	return nil
}

func resourceLaceworkTeamMemberRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceLaceworkTeamMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceLaceworkTeamMemberRead(d, meta)
}

func resourceLaceworkTeamMemberDelete(d *schema.ResourceData, meta interface{}) error {
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

