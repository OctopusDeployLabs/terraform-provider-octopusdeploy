package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandKubernetesCluster(d *schema.ResourceData) *octopusdeploy.KubernetesEndpoint {
	clusterURL, _ := url.Parse(d.Get("cluster_url").(string))

	endpoint := octopusdeploy.NewKubernetesEndpoint(clusterURL)
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("authentication"); ok {
		endpoint.Authentication = expandEndpointAuthentication(v)
	}

	if v, ok := d.GetOk("cluster_certificate"); ok {
		endpoint.ClusterCertificate = v.(string)
	}

	if v, ok := d.GetOk("container"); ok {
		endpoint.Container = expandDeploymentActionContainer(v)
	}

	if v, ok := d.GetOk("default_worker_pool_id"); ok {
		endpoint.DefaultWorkerPoolID = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		endpoint.Namespace = v.(string)
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		endpoint.ProxyID = v.(string)
	}

	if v, ok := d.GetOk("running_in_container"); ok {
		endpoint.RunningInContainer = v.(bool)
	}

	if v, ok := d.GetOk("skip_tls_verification"); ok {
		endpoint.SkipTLSVerification = v.(bool)
	}

	return endpoint
}

func flattenKubernetesCluster(endpoint *octopusdeploy.KubernetesEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"authentication":         flattenEndpointAuthentication(endpoint.Authentication),
		"cluster_certificate":    endpoint.ClusterCertificate,
		"container":              flattenDeploymentActionContainer(endpoint.Container),
		"default_worker_pool_id": endpoint.DefaultWorkerPoolID,
		"id":                     endpoint.GetID(),
		"namespace":              endpoint.Namespace,
		"proxy_id":               endpoint.ProxyID,
		"running_in_container":   endpoint.RunningInContainer,
		"skip_tls_verification":  endpoint.SkipTLSVerification,
	}

	if endpoint.ClusterURL != nil {
		rawEndpoint["cluster_url"] = endpoint.ClusterURL.String()
	}

	return []interface{}{rawEndpoint}
}

func getKubernetesClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"authentication": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getEndpointAuthenticationSchema()},
			MaxItems: 1,
			MinItems: 0,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"cluster_certificate": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"cluster_url": {
			Required: true,
			Type:     schema.TypeString,
		},
		"container": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getDeploymentActionContainerSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"default_worker_pool_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"namespace": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"proxy_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"running_in_container": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"skip_tls_verification": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
