# taskapp

`taskapp` is a simple task list app. Users can create, modify, list, and delete tasks. Daily task list summaries are e-mailed to users. Users can also create notes and collaborate with other users on them.

I used [apex](http://apex.run) for AWS lambda packaging/deployment. I used [Terraform](https://www.terraform.io/) for configuration and deployment of DynamoDB, API Gateway, and Cloudwatch alarms.

## Project layout

```
* functions - Lambda function implementations, all in Go
* infrastructure - Terraform (.tf) files for deployment of DynamoDB, API Gateway
* shared - Go package with shared utilities for the functions
* test - Sanity checks for the system
* bootstrap.sh - Installs apex and terraform on host machine.
* deploy.sh - Deploys everything using apex, terraform
* deps.sh - Installs go dependencies, of which there are few
* init.sh - The script I ran to initiate the apex project (for completeness)
* project.json - The apex project config
* lambda_iam.json - The AWS IAM role for the lambdas.
* README.md - See README.md
* swagger.json - A swagger definition for the complete API
```

## Development

### Initialization

I originally initialized the project with apex using the following script:

```
$ ./init.sh
```

There's no need to call this again.

### Bootstrap

For development of this project, you will need [apex](http://apex.run) and [Terraform](https://www.terraform.io/) installed. The following script will install both of those on your machine.

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

#### Lambda IAM Role

Apex requires the IAM roles for the Lambdas to be created before deployment. This is not ideal. It would be preferrable to generate them with terraform. Presently it's necessary to manually generate the Lambda IAM role and policy and place the IAM role ARN in the project.json directory.

#### Apex credentials

For apex, I stored my AWS config and credentials in `~/.aws/config` and `~/.aws/credentials` using the `taskapp` profile.

#### Terraform credentials

For terraform, I stored my AWS credentials in `infrastructure/dev/secret.tfvars`. The file looks like this:

```
account_id = "XXXXXXXXXXXX"
access_key = "XXXXXXXXXXXXXXXXXXXX"
secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

### Deploy

To deploy the Lambdas (using Apex), DynamoDB, API Gateway, and Cloudwatch alarms (using Terraform) simply call:

```
$ ./deploy.sh
```

This will also deploy a stage of the API called `dev`. There are more details in the `infrastructure/dev` directory. 

## Features

### API

There is a swagger.json file included that provides a relatively complete documentation of the API.

### Tests

The tests for this project are more sanity checks than actual tests. They are simple curl scripts necessary to invoke the API from its deployed state. For example, to add a task via the deployed dev stage, do:

```
test/apigateway/task_add0.sh
```

There are various similar tests in the directory. These are provided as useful debugging tools and could become proper automated integration tests in the future.


[Example list call](https://ovfepswc3l.execute-api.us-east-1.amazonaws.com/dev/task/list?user=testy.mctester@example.com)

* Incorporated unique, randomly generated ids into the project. Greatly simplified the API's.
* You'll need to manually set the AWS IAM credentials in project.json
* Example: https://ovfepswc3l.execute-api.us-east-1.amazonaws.com/dev/task/list?user=bar@baz.com
* Could add many more error codes but ran out of time
* Sort of /task/list is O(n*logn)
* Could have more robust tests, a load generator, better logging, but ran out of time.
* Apex and terraform integrate poorly. It is painful to use IAM roles generated from terraform with apex.
* Wish that Terraform supported swagger.json input.