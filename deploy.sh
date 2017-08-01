#/bin/bash

# use apex to deploy all of the lambdas
apex deploy

# use terraform to deploy dynamodb, api gateway
apex infra apply -var-file=secret.tfvars