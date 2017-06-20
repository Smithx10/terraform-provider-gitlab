package gitlab

import (
	"fmt"
	gitlab "go-gitlab"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

// bob
func resourceGitlabGroupMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupMemberCreate,
		Read:   resourceGitlabGroupMemberRead,
		Update: resourceGitlabGroupMemberUpdate,
		Delete: resourceGitlabGroupDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"accesslevel": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
	}
}

func resourceGitlabGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.AddGroupMemberOptions{
		UserID: gitlab.Int(d.Get("user_id").(int)),
	}

	log.Printf("[DEBUG] create gitlab member %q", options.UserID)

	member, _, err := client.Groups.AddGroupMember("test-group", options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", member.ID))

	return nil
}

func resourceGitlabGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGitlabGroupMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGitlabGroupMemberDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
