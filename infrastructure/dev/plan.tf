provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

// define dynamodb table
// define iam policies
// inject iam policies into functions
// (apex inject lambdas)
// define api gateway referencing the lambdas

// DynamoDB table

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




/*

resource "aws_api_gateway_rest_api" "taskapp" {
  name        = "taskapp"
  description = "This is my API for demonstration purposes"
}

resource "aws_api_gateway_resource" "MyDemoResource" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_rest_api.taskapp.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_method" "taskapp_add" {
  rest_api_id   = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id   = "${aws_api_gateway_resource.MyDemoResource.id}"
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "MyDemoIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.MyDemoResource.id}"
  http_method = "${aws_api_gateway_method.taskapp_add.http_method}"
  type        = "LAMBDA"
}

resource "aws_api_gateway_deployment" "MyDemoDeployment" {
  depends_on = ["aws_api_gateway_method.taskapp_add"]

  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  stage_name  = "test"

  variables = {
    "answer" = "42"
  }
}
*/