provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}


# Until now, the resource created could not respond to anything. We must set up
# a HTTP method (or verb) for that!
# This is the code for method GET /hello, that will talk to the first lambda
module "hello_get" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  resource_id = "${aws_api_gateway_resource.example_api_res_hello.id}"
  method      = "GET"
  path        = "${aws_api_gateway_resource.example_api_res_hello.path}"
  lambda      = "${var.apex_function_list}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
}

# This is the code for method POST /hello, that will talk to the second lambda
module "hello_post" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  resource_id = "${aws_api_gateway_resource.example_api_res_hello.id}"
  method      = "POST"
  path        = "${aws_api_gateway_resource.example_api_res_hello.path}"
  lambda      = "${var.apex_function_add}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
}



# DynamoDB
resource "aws_dynamodb_table" "taskapp_table" {
  name           = "taskapp_terr"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  provisioner "local-exec" {
    command = "echo ${self.name} >> tablename.txt"
  }
}

# Now, we need an API to expose those functions publicly
resource "aws_api_gateway_rest_api" "example_api" {
  name = "Hello API"
}

# The API requires at least one "endpoint", or "resource" in AWS terminology.
# The endpoint created here is: /hello
resource "aws_api_gateway_resource" "example_api_res_hello" {
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  parent_id   = "${aws_api_gateway_rest_api.example_api.root_resource_id}"
  path_part   = "hello"
}

resource "aws_api_gateway_deployment" "example_api_deployment" {
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  stage_name  = "dev"
  description = "Deploy methods: ${module.hello_get.http_method} ${module.hello_post.http_method}"
}


/*


# API Gateway
resource "aws_api_gateway_rest_api" "example_api" {
  name = "ExampleAPI"
  description = "Example Rest Api"
}

# API Gateway - /task/list
resource "aws_api_gateway_resource" "example_api_resource" {
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  parent_id = "${aws_api_gateway_rest_api.example_api.root_resource_id}"
  path_part = "list"
}

resource "aws_api_gateway_method" "example_api_method" {
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  resource_id = "${aws_api_gateway_resource.example_api_resource.id}"
  http_method = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "example_api_method-integration" {
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  resource_id = "${aws_api_gateway_resource.example_api_resource.id}"
  http_method = "${aws_api_gateway_method.example_api_method.http_method}"
  type = "AWS"
  integration_http_method = "POST"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.apex_function_add}/invocations"
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${var.apex_function_add}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_api_gateway_rest_api.example_api.id}/&@/${aws_api_gateway_method.example_api_method.http_method}/resourcepath/subresourcepath"
}

# API Gateway - Deployment

resource "aws_api_gateway_deployment" "example_deployment_dev" {
  depends_on = [
    "aws_api_gateway_method.example_api_method",
    "aws_api_gateway_integration.example_api_method-integration"
  ]
  rest_api_id = "${aws_api_gateway_rest_api.example_api.id}"
  stage_name = "dev"
}

output "dev_url" {
  value = "https://${aws_api_gateway_deployment.example_deployment_dev.rest_api_id}.execute-api.${var.region}.amazonaws.com/${aws_api_gateway_deployment.example_deployment_dev.stage_name}"
}

*/

/*

# Due to limitations with apex, it doesn't seem possible to use deploy specific IAM policies 

resource "aws_iam_role" "test_role" {
  name = "test_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  provisioner "local-exec" {
    command = "echo ${self.unique_id} >> iam_role_id.txt"
  }
}

resource "aws_iam_role_policy" "test_policy" {
  name = "test_policy"
  role = "${aws_iam_role.test_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "dynamodb:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

*/