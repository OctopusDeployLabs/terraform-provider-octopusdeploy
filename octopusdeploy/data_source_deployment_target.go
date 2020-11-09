package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeploymentTargetReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"endpoint": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"communication_style": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"proxy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"thumbprint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster_certificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cluster_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"skip_tls_verification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"default_worker_pool_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 0,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"client_certificate": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"authentication_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
											"KubernetesCertificate",
											"KubernetesStandard",
										}, false)),
									},
								},
							},
						},
					},
				},
			},
			"environments": {
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
			"has_latest_calamari": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_disabled": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_in_process": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"machine_policy_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"roles": {
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status_summary": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tenanted_deployment_participation": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tenant_ids": {
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
			"tenant_tags": {
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceDeploymentTargetReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	deploymentTargets, err := client.Machines.GetByName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(deploymentTargets) == 0 {
		return diag.Errorf("unable to retrieve deployment target (filter: %s)", name)
	}

	// NOTE: two or more deployment targets can have the same name in Octopus
	// and therefore, a better search criteria needs to be implemented below

	for _, deploymentTarget := range deploymentTargets {
		if deploymentTarget.Name == name {
			flattenDeploymentTarget(ctx, d, deploymentTarget)
			return nil
		}
	}

	return nil
}
