# Running integration tests

At this moment you can run tests in two ways:
- Using Octopus Deploy container created within test session
- Using separately running instance of Octopus Deploy (BYO)

> [!WARNING]
> We have tests which directly access automatically created Octopus Deploy Docker container and executes 'terraform' command directly.  
> For such tests make sure that `terraform` command is recognisable.  
> We are planning to remove these and replace with tests which rely on terraform built-in testing framework

## Separately running Octopus Deploy instance (BYO)

Run instance of Octopus Deploy before executing tests.

Make sure next environment variables are available:
```
TF_ACC_LOCAL=1
OCTOPUS_URL='http://localhost:<port-number>'
OCTOPUS_APIKEY='<api-key>'
```

### From terminal
Execute from repository's root directory
```
go test -run "^(?:TestAccOctopusDeployNuGetFeedBasic|TestAccResourceBuiltInTrigger)$" -timeout 0 ./...
```

## Automatically created Octopus Deploy instance

Instance of Octopus Deploy and it's dependencies will be created and run automatically as part of the test session.  

Make sure next environment variables are available:

```
LICENSE=<octopus-deploy-server-license-base-64-string>
OCTOTESTIMAGEURL='docker.packages.octopushq.com/octopusdeploy/octopusdeploy'
OCTOTESTVERSION='latest'
OCTODISABLEOCTOCONTAINERLOGGING='true'
OCTOTESTSKIPINIT='true'
OCTOTESTRETRYCOUNT=1
GOMAXPROCS=1
```

### From terminal
Execute from repository's root directory with additional parameter `-createSharedContainer=true`
```
go test -run "^(?:TestAccResourceBuiltInTrigger)$" -timeout 0 ./... -createSharedContainer=true
```

## Testing Environment
Test may require elevated privileges to create symlinks for schema directories

### Optional environment variables
Some tests require next environment variables    
AWS related:
```
ECR_ACCESS_KEY
ECR_SECRET_KEY
```
GitHub related:
```
GIT_USERNAME
GIT_CREDENTIAL
```

### Configure IDE (GoLand)
GoLand creates new run configuration for every new test.  
To be able to run different tests without manually adding environment variables to new configuration, we can set them once in "Go Test" configuration template.   
`Edit Configurations...` >> `Edit configuration templates...` >> `Go Test`

- Add environment variables
- Enable `Run with elevated privileges` (_IDE creating symlinks for schema directories_)

Now you will be able to run or debug tests

