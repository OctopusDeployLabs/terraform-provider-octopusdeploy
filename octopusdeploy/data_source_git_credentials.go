package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGitCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing GitCredentials.",
		ReadContext: dataSourceGitCredentialsRead,
		Schema:      getGitCredentialDataSchema(),
	}
}

func dataSourceGitCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := credentials.Query{
		Name: d.Get("name").(string),
		Skip: d.Get("skip").(int),
		Take: d.Get("take").(int),
	}
	spaceID := d.Get("space_id").(string)

	client := m.(*client.Client)
	existingGitCredentials, err := credentials.Get(client, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedGitCredentials := []interface{}{}
	for _, gitCredential := range existingGitCredentials.Items {
		flattenedGitCredentials = append(flattenedGitCredentials, flattenGitCredential(gitCredential))
	}

	d.Set("git_credentials", flattenedGitCredentials)
	d.SetId("GitCredentials " + time.Now().UTC().String())

	return nil
}
