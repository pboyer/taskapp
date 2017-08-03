# taskapp

`taskapp` is a simple task list app. I used [apex](http://apex.run) for AWS lambda packaging/deployment. I used [Terraform](https://www.terraform.io/) for deployment of DynamoDB and API Gateway.

* I stored my aws config and credentials in ~/.aws/config and ~/.aws/credentials
* You'll need to manually set the lambda credentials in project.json
* Secrets should be stored in infrastructure/dev/secret.tfvars
* Completed field contains null when there is no data
* Example: https://ovfepswc3l.execute-api.us-east-1.amazonaws.com/dev/task/list?user=bar@baz.com

## Project layout

```
* functions - Lambda function implementations, all in Go
* infrastructure - Terraform (.tf) files for deployment of DynamoBB, API Gateway
* shared - Go package with shared utilities for the functions
* test - Sanity tests for the system
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

