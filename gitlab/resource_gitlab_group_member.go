package gitlab

import (
	"fmt"
	"log"

	gitlab "github.com/xanzy/go-gitlab"
	"github.com/y0ssar1an/q"

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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
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

func resourceGitlabGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	// Bring in the gitlab Client
	client := meta.(*gitlab.Client)

	// Convert user_name to UserID
	u, err := convertUNtoID(d.Get("user_name").(string), meta)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// Type Covnert schema.TypeInt to gitlab.AccessLevelValue
	a, ok := d.Get("access_level").(gitlab.AccessLevelValue)
	if ok {
		fmt.Printf("Int value is %d\n", a)
	} else {
		fmt.Println("wrong type for access level, expected int")
	}
	q.Q(a, ok, d.Get("access_level"))
	// Define AddGroupMemberOptions with values from schema
	l := &gitlab.AddGroupMemberOptions{
		UserID:      gitlab.Int(u),
		AccessLevel: gitlab.AccessLevel(a),
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
