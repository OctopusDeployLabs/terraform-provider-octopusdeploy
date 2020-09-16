# Resource Certificate Resource

## Required Values
The required values that need to be passed into the Terraform config block for this resource are:
1. name - The name of the newly created resource
2. certificate_data - The actual certificate data (the certificate itself) that's being passed into Octopus Deploy

## Optional Values
The optional values that can be apssed into the Terraform config block for this resource are:
1. notes- The description for the newly created resource
2. enviromments - Environments to target for this resource
3. tenant_tags - Any tags for this resource

## What can this resource do?
The resource can perform the following actions:

1. Create
2. Update
3. Delete