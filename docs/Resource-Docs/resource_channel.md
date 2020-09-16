# Channel Resource

## Required Values
The required values that need to be passed into the Terraform config block for this resource are:
1. name - The name of the newly created resource
2. project_id - The project ID that the channel is being created in
3. key - The Azure app registration secret/password key
4. subscription_number - The Azure subscription that you will be creating the resource in

## Optional Values
The optional values that can be apssed into the Terraform config block for this resource are:
1. description - The description for the newly created resource
2. lifecycle_id - The lifecycle ID that the channel will be associated with
3. is_default - Specifying whether or not the channel is the default. The default value for this parameter is `true`
4. rule - A rule for the channel that can specify the version range, tag, and actions

## What can this resource do?
The resource can perform the following actions:

1. Create
2. Update
3. Delete