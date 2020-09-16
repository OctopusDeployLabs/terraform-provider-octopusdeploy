# Azure Service Principal Account Resource

## Required Values
The required values that need to be passed into the Terraform config block for this resource are:
1. name - The name of the newly created resource
2. client_id - The Azure app registration client/app ID
3. key - The Azure app registration secret/password key
4. subscription_number - The Azure subscription that you will be creating the resource in

## Optional Values
The optional values that can be apssed into the Terraform config block for this resource are:
1. description - The description for the newly created resource
2. enviromments - Environments to target for this resource
3. tenant_tags - Any tags for this resource
4. resource_management_endpoint_base_uri - URI for the resource manager endpoint if one is used
5. active_directory_endpoint_base_uri - Active Directory URI if this app registration is managed by Active Directory

## What can this resource do?
The resource can perform the following actions:

1. Create
2. Update
3. Delete