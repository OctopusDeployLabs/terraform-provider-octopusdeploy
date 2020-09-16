# AWS Account Resource

## Required Values
The required values that need to be passed into the Terraform config block for this resource are:
1. name - The name of the newly created resource
2. access_key - The AWS IAM access key
3. secret_key - The AWS IAM secret key

## Optional Values
The optional values that can be apssed into the Terraform config block for this resource are:
1. description - The description for the newly created resource
2. enviromments - Environments to target for this resource
3. tenant_tags - Any tags for this resource

## What can this resource do?
The resource can perform the following actions:

1. Create
2. Update
3. Delete