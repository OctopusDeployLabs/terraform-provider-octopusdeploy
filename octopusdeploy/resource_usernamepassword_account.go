package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUsernamePassword() *schema.Resource {

	schemaMap := getCommonAccountsSchema()

	schemaMap["username"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	schemaMap["password"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return &schema.Resource{
		Create: resourceUsernamePasswordCreate,
		Read:   resourceUsernamePasswordRead,
		Update: resourceUsernamePasswordUpdate,
		Delete: resourceAccountDeleteCommon,
		Schema: schemaMap,
	}
}

func resourceUsernamePasswordRead(d *schema.ResourceData, m interface{}) error {
	_, err := fetchAndReadAccount(d, m);
	//account, err := fetchAndReadAccount(d, m);

	if err != nil {
		return err;
	}

	// TODO: Username property does not yet exist
	//d.Set("username", account.Username)

	return nil;
}

func buildUsernamePasswordResource(d *schema.ResourceData) *octopusdeploy.Account {
	account := buildAccountResourceCommon(d, octopusdeploy.UsernamePassword);

	// TODO: Username property does not yet exist
	//if v, ok := d.GetOk("username"); ok {
	//	account.Username = v.(string)
	//}

	if v, ok := d.GetOk("password"); ok {
		account.Password = octopusdeploy.SensitiveValue{NewValue: v.(string)}
	}

	return account
}

func resourceUsernamePasswordCreate(d *schema.ResourceData, m interface{}) error {
	account := buildUsernamePasswordResource(d)
	return resourceAccountCreateCommon(d, m, account);
}

func resourceUsernamePasswordUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildUsernamePasswordResource(d)
	return resourceAccountUpdateCommon(d, m, account)
}
