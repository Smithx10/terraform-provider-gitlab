package gitlab

import (
	"fmt"
	"log"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGitlabGroupMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupMemberCreate,
		Read:   resourceGitlabGroupMemberRead,
		Update: resourceGitlabGroupMemberUpdate,
		Delete: resourceGitlabGroupMemberDelete,

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateValueFunc([]string{"GuestPermissions", "ReporterPermissions", "DeveloperPermissions", "MasterPermissions", "OwnerPermission"}),
				Default:      "DeveloperPermissions",
			},
		},
	}
}

// function to convert usernames to userids
func convertUNtoID(un string, meta interface{}) (int, error) {
	client := meta.(*gitlab.Client)

	l := &gitlab.ListUsersOptions{
		Username: gitlab.String(un),
	}

	u, _, err := client.Users.ListUsers(l)
	if err != nil {
		return fmt.Printf("%v", err)
	}
	if un != u[0].Username {
		return fmt.Printf("%v", un)
	}

	return u[0].ID, nil
}

// Handle String to AccessLevelValue Conversion
func getAccessLevel(al string) gitlab.AccessLevelValue {
	m := make(map[string]gitlab.AccessLevelValue)
	m["GuestPermissions"] = gitlab.GuestPermissions
	m["ReporterPermissions"] = gitlab.ReporterPermissions
	m["DeveloperPermissions"] = gitlab.DeveloperPermissions
	m["DeveloperPermissions"] = gitlab.MasterPermissions
	m["OwnerPermission"] = gitlab.OwnerPermission
	return m[al]
}

func resourceGitlabGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	// Bring in the gitlab Client
	client := meta.(*gitlab.Client)

	// Convert Schema user_name.(string) to gitlab.ID.(int)
	u, err := convertUNtoID(d.Get("user_name").(string), meta)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// Convert User input access_level.(string) to gitlab.AccessLevelValue
	al := getAccessLevel(d.Get("access_level").(string))

	// Define AddGroupMemberOptions with values from schema
	l := &gitlab.AddGroupMemberOptions{
		UserID:      gitlab.Int(u),
		AccessLevel: gitlab.AccessLevel(al),
	}

	log.Printf("[DEBUG] create gitlab member %q", l.UserID)

	// Execute the Addgroup MemberShip
	m, _, err := client.Groups.AddGroupMember(d.Get("group_name"), l)
	if err != nil {
		return err
	}

	// SetId to assert we created the record
	d.SetId(fmt.Sprintf("%d", m.ID))

	return nil
}

func resourceGitlabGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	fmt.Print("insideread")
	return nil
}

func resourceGitlabGroupMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	fmt.Print("insideupdate")
	return nil
}

func resourceGitlabGroupMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	t := d.Get("group_name")

	// Aquire the UserID based on Username
	n, err := convertUNtoID(d.Get("user_name").(string), meta)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	_, err1 := client.Groups.RemoveGroupMember(t, n)
	return err1
}
