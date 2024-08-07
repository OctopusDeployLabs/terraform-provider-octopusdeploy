package octopusdeploy

import (
	"context"
	"log"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func resourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		DeleteContext: resourceTagDelete,
		Description:   "This resource manages tags in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTagRead,
		Schema:        getTagSchema(),
		UpdateContext: resourceTagUpdate,
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	log.Printf("[INFO] creating tag")

	return tagCreate(ctx, d, m)
}

func tagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] creating tag")

	tagSetID := d.Get("tag_set_id").(string)
	tagSetSpaceID := d.Get("tag_set_space_id").(string)

	octopus := m.(*client.Client)
	tagSet, err := tagsets.GetByID(octopus, tagSetSpaceID, tagSetID)
	if err != nil {
		return processUnknownTagSetError(ctx, d, err)
	}

	name := d.Get("name").(string)

	for _, tag := range tagSet.Tags {
		if tag.Name == name {
			return diag.Errorf(`the tag name '%s' is already in use by another tag in this tag set; tag names must be unique`, name)
		}
	}

	tag := expandTag(d)
	if tag.ID != "" {
		tag.ID = tagSet.GetID() + "/" + strings.Split(tag.ID, "/")[1]
	}
	tagSet.Tags = append(tagSet.Tags, tag)

	updatedTagSet, err := tagsets.Update(octopus, tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	return findByIdOrNameAndSetTag(ctx, d, tag, updatedTagSet)

}

func resourceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	tagSetID := d.Get("tag_set_id").(string)
	tagSetSpaceID := d.Get("tag_set_space_id").(string)

	log.Printf("[INFO] deleting tag (%s)", d.Id())

	octopus := m.(*client.Client)
	tagSet, err := tagsets.GetByID(octopus, tagSetSpaceID, tagSetID)
	if err != nil {
		return processUnknownTagSetError(ctx, d, err)
	}

	tag := expandTag(d)

	// verify tag is not associated with a tenant
	isUsed, err := isTagUsedByTenants(ctx, octopus, tagSetSpaceID, tag)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if isUsed {
		d.SetId("")
		return diag.Errorf("the tag may not be deleted; it is being used by one or more tenant(s)")
	}

	// tag is known and not associated with a tenant, therefore it may be deleted

	for i := 0; i < len(tagSet.Tags); i++ {
		if tagSet.Tags[i].ID == d.Id() {
			tagSet.Tags = slices.Delete(tagSet.Tags, i, i+1)

			if _, err := tagsets.Update(octopus, tagSet); err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[INFO] tag deleted (%s)", d.Id())
			d.SetId("")
			return nil
		}
	}

	return errors.DeleteFromState(ctx, d, "tag")
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	// validate the tag ID
	if d.Id() == "" || !strings.Contains(d.Id(), "/") {
		return diag.Errorf(`unable to import tag; ID must be "TagSets-{ID}/Tags-{ID}"`)
	}

	name := d.Get("name").(string)
	tagSetID := d.Get("tag_set_id").(string)
	tagSetSpaceID := d.Get("tag_set_space_id").(string)

	// if name and tag set ID are empty then an import is underway
	if name == "" && tagSetID == "" {
		log.Printf("[INFO] importing tag (%s)", d.Id())
		tagSetID = strings.Split(d.Id(), "/")[0]
	} else {
		log.Printf("[INFO] reading tag (%s)", d.Id())
	}

	octopus := m.(*client.Client)
	tagSet, err := tagsets.GetByID(octopus, tagSetSpaceID, tagSetID)
	if err != nil {
		return processUnknownTagSetError(ctx, d, err)
	}

	tag := expandTag(d)
	return findByIdOrNameAndSetTag(ctx, d, tag, tagSet)
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	name := d.Get("name").(string)
	tagSetID := d.Get("tag_set_id").(string)
	tagSetSpaceID := d.Get("tag_set_space_id").(string)

	log.Printf("[INFO] updating tag (%s)", d.Id())

	octopus := m.(*client.Client)

	// if the tag is reassigned to another tag set
	if d.HasChange("tag_set_id") {
		sourceTagSetID, destinationTagSetID := d.GetChange("tag_set_id")
		sourceTagSetSpaceID, destinationTagSetSpaceID := d.GetChange("tag_set_space_id")

		sourceTagSet, err := tagsets.GetByID(octopus, sourceTagSetSpaceID.(string), sourceTagSetID.(string))
		if err != nil {
			// if spaceID has change, tag has been deleted, recreate required
			if d.HasChange("tag_set_space_id") {
				return tagCreate(ctx, d, m)
			}
			return diag.FromErr(err)
		}

		destinationTagSet, err := tagsets.GetByID(octopus, destinationTagSetSpaceID.(string), destinationTagSetID.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		// check to see if the name already exists in the destination tag set
		for _, tag := range destinationTagSet.Tags {
			if tag.Name == name {
				d.SetId("")
				return diag.Errorf(`the tag name '%s' is already in use by another tag in this tag set; tag names must be unique`, name)
			}
		}

		tag := expandTag(d)

		// check to see that the tag is not applied to a tenant
		isUsed, err := isTagUsedByTenants(ctx, octopus, sourceTagSetSpaceID.(string), tag)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		if isUsed {
			d.SetId("")
			return diag.Errorf("the tag may not be transferred; it is being used by one or more tenant(s)")
		}

		// all requirements are met; it is OK to transfer the tag

		// remove the tag from the source tag set and update through the API
		for i := 0; i < len(sourceTagSet.Tags); i++ {
			if sourceTagSet.Tags[i].ID == d.Id() {
				sourceTagSet.Tags = slices.Delete(sourceTagSet.Tags, i, i+1)

				if _, err := tagsets.Update(octopus, sourceTagSet); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// update and add the tag to the destination tag set
		tag.ID = destinationTagSet.GetID() + "/" + strings.Split(tag.ID, "/")[1]
		destinationTagSet.Tags = append(destinationTagSet.Tags, tag)

		updatedTagSet, err := tagsets.Update(octopus, destinationTagSet)
		if err != nil {
			return diag.FromErr(err)
		}

		return findByIdOrNameAndSetTag(ctx, d, tag, updatedTagSet)
	}

	tagSet, err := tagsets.GetByID(octopus, tagSetSpaceID, tagSetID)
	if err != nil {
		return processUnknownTagSetError(ctx, d, err)
	}

	// find and update the tag that matches the one updated in configuration
	for i := 0; i < len(tagSet.Tags); i++ {
		if tagSet.Tags[i].ID == d.Id() {
			tagSet.Tags[i] = expandTag(d)

			updatedTagSet, err := tagsets.Update(octopus, tagSet)
			if err != nil {
				return diag.FromErr(err)
			}

			return findByIdOrNameAndSetTag(ctx, d, tagSet.Tags[i], updatedTagSet)
		}
	}

	return diag.Errorf("unable to update tag")
}

func isTagUsedByTenants(ctx context.Context, octopus *client.Client, spaceID string, tag *tagsets.Tag) (bool, error) {
	tenants, err := tenants.Get(octopus, spaceID, tenants.TenantsQuery{
		Tags: []string{tag.ID},
	})
	if err != nil {
		return false, err
	}

	return len(tenants.Items) > 0, nil
}

func findByIdOrNameAndSetTag(ctx context.Context, d *schema.ResourceData, tag *tagsets.Tag, tagSet *tagsets.TagSet) diag.Diagnostics {
	for _, t := range tagSet.Tags {
		if t.Name == tag.Name {
			if err := setTag(ctx, d, t, tagSet); err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[INFO] tag (%s)", tag.ID)
			return nil
		}
	}

	for _, t := range tagSet.Tags {
		if t.ID == tag.ID {
			if err := setTag(ctx, d, t, tagSet); err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] tag (%s)", tag.ID)
			return nil
		}
	}

	return errors.DeleteFromState(ctx, d, "tag")
}

func processUnknownTagSetError(ctx context.Context, d *schema.ResourceData, err error) diag.Diagnostics {
	if err == nil {
		return nil
	}

	if apiError, ok := err.(*core.APIError); ok {
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] tag set (%s) not found; deleting tag from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	return diag.FromErr(err)
}
