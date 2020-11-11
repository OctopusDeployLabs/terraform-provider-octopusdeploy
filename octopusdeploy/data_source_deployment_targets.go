package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDeploymentTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeploymentTargetsRead,
		Schema: map[string]*schema.Schema{
			"communication_styles": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"deployment_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"environments": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"health_statuses": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"ids": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"is_disabled": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"partial_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"roles": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"shell_names": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"skip": {
				Default:  0,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"take": {
				Default:  1,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tenants": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"tenant_tags": {
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Type:     schema.TypeList,
			},
			"thumbprint": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"deployment_targets": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication": {
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"account_id": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"client_certificate": {
													Optional: true,
													Type:     schema.TypeString,
												},
												"authentication_type": {
													Optional: true,
													Type:     schema.TypeString,
													ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
														"KubernetesCertificate",
														"KubernetesStandard",
													}, false)),
												},
											},
										},
										Optional: true,
										Type:     schema.TypeSet,
									},
									"cluster_certificate": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cluster_url": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"communication_style": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"default_worker_pool_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"namespace": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"proxy_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"skip_tls_verification": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"thumbprint": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Type: schema.TypeList,
						},
						"environments": {
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Type:     schema.TypeList,
						},
						"has_latest_calamari": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"health_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"id": {
							Computed: true,
							Type:     schema.TypeString,
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
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"operating_system": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"roles": {
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Type:     schema.TypeList,
						},
						"shell_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"shell_version": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"space_id": {
							Computed: true,
							Type:     schema.TypeString,
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
						"tenants": {
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Type:     schema.TypeList,
						},
						"tenant_tags": {
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Type:     schema.TypeList,
						},
						"thumbprint": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"uri": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceDeploymentTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.MachinesQuery{
		CommunicationStyles: expandArray(d.Get("communication_styles").([]interface{})),
		DeploymentID:        d.Get("deployment_id").(string),
		EnvironmentIDs:      expandArray(d.Get("environments").([]interface{})),
		HealthStatuses:      expandArray(d.Get("health_statuses").([]interface{})),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsDisabled:          d.Get("is_disabled").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		Roles:               expandArray(d.Get("roles").([]interface{})),
		ShellNames:          expandArray(d.Get("shell_names").([]interface{})),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
		TenantIDs:           expandArray(d.Get("tenants").([]interface{})),
		TenantTags:          expandArray(d.Get("tenant_tags").([]interface{})),
		Thumbprint:          d.Get("thumbprint").(string),
	}

	client := m.(*octopusdeploy.Client)
	deploymentTargets, err := client.Machines.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedDeploymentTargets := []interface{}{}
	for _, deploymentTarget := range deploymentTargets.Items {
		flattenedDeploymentTarget := map[string]interface{}{
			"endpoint":                          flattenEndpoint(deploymentTarget.Endpoint),
			"environments":                      deploymentTarget.EnvironmentIDs,
			"has_latest_calamari":               deploymentTarget.HasLatestCalamari,
			"health_status":                     deploymentTarget.HealthStatus,
			"id":                                deploymentTarget.GetID(),
			"is_disabled":                       deploymentTarget.IsDisabled,
			"is_in_process":                     deploymentTarget.IsInProcess,
			"machine_policy_id":                 deploymentTarget.MachinePolicyID,
			"name":                              deploymentTarget.Name,
			"operating_system":                  deploymentTarget.OperatingSystem,
			"roles":                             deploymentTarget.Roles,
			"shell_name":                        deploymentTarget.ShellName,
			"shell_version":                     deploymentTarget.ShellVersion,
			"space_id":                          deploymentTarget.SpaceID,
			"status":                            deploymentTarget.Status,
			"status_summary":                    deploymentTarget.StatusSummary,
			"tenanted_deployment_participation": deploymentTarget.TenantedDeploymentMode,
			"tenants":                           deploymentTarget.TenantIDs,
			"tenant_tags":                       deploymentTarget.TenantTags,
			"thumbprint":                        deploymentTarget.Thumbprint,
			"uri":                               deploymentTarget.URI,
		}
		flattenedDeploymentTargets = append(flattenedDeploymentTargets, flattenedDeploymentTarget)
	}

	d.Set("deployment_targets", flattenedDeploymentTargets)
	d.SetId("DeploymentTargets " + time.Now().UTC().String())

	return nil
}
