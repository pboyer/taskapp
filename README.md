# taskapp

`taskapp` is a simple task list app. I used [apex](http://apex.run) for AWS lambda packaging/deployment. I used [Terraform](https://www.terraform.io/) for deployment of DynamoDB, API Gateway, and Cloudwatch alarms.

## Repo layout

```
* functions - Lambda function implementations, all in Go
* infrastructure - Terraform (.tf) files for deployment of DynamoDB, API Gateway
* shared - Go package with shared utilities for the functions
* test - Sanity checks for the system
* bootstrap.sh - Installs apex, terraform on host machine
* deploy.sh - Deploys everything using apex, terraform
* deps.sh - Installs go dependencies, of which there are few
* init.sh - The script I ran to initiate the apex project (for completeness)
* project.json - The apex project config
* README.md - See README.md
* swagger.json - A swagger definition for the complete API
```

## Development

### Initialization

I originally initialized the project with apex using the following script:

```
$ ./init.sh
```

### Bootstrap

For development of this project, you will need [apex](http://apex.run) and [Terraform](https://www.terraform.io/) installed. The following script will install both of those.

```
$ ./bootstrap.sh
```

You will also need Go installed and have this project on your `GOPATH`. This project has only been tested with:

```
go version go1.7.3 darwin/amd64
```

Before development, you'll need to call to obtain all of the Go packages in use:

```
$ ./bootstrap.sh
```

### Deploy

To deploy the lambdas, DynamoDB, and API Gateway simply call:

```
$ ./deploy.sh
```

[Example list call](https://ovfepswc3l.execute-api.us-east-1.amazonaws.com/dev/task/list?user=testy.mctester@example.com)

* Incorporated unique, randomly generated ids into the project. Greatly simplified the API's.
* I stored my aws config and credentials in ~/.aws/config and ~/.aws/credentials
* You'll need to manually set the lambda credentials in project.json
* Secrets should be stored in infrastructure/dev/secret.tfvars
* Completed field contains null when there is no data
* Example: https://ovfepswc3l.execute-api.us-east-1.amazonaws.com/dev/task/list?user=bar@baz.com
* Could add many more error codes but ran out of time
* Sort of /task/list is O(n*logn)
* Could have more robust tests, a load generator, better logging, but ran out of time.
* Apex and terraform integrate somewhat poorly. Painful to use IAM roles generated from terraform with apex.
* Wish that Terraform supported swagger.json input