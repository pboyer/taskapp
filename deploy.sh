#/bin/bash

# use apex to deploy all of the lambdas
apex deploy

# Get terraform modules
apex infra get

# use terraform to deploy dynamodb, api gateway
apex infra apply -var-file=secret.tfvars