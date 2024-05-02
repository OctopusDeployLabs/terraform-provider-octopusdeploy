package octopusdeploy

//func expandKubernetesTentacleEndpoint(flattenedMap map[string]interface{}) *machines.KubernetesTentacleEndpoint {
//	thumbprint := flattenedMap["thumbprint"].(string)
//	uri, _ := url.Parse(flattenedMap["uri"].(string))
//	communicationsMode := flattenedMap["communication_mode"].(string)
//	defaultNamespace := flattenedMap["namespace"].(string)
//	upgradeLocked := flattenedMap["upgrade_locked"].(bool)
//	endpoint := machines.NewKubernetesTentacleEndpoint(uri, thumbprint, upgradeLocked, communicationsMode, defaultNamespace)
//	endpoint.ID = flattenedMap["id"].(string)
//
//	return endpoint
//}
