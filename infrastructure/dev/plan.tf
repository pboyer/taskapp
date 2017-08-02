provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

# DynamoDB

resource "aws_dynamodb_table" "taskapp" {
  name           = "taskapp"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

# API Gateway

resource "aws_api_gateway_rest_api" "taskapp" {
  name = "Task App"
}

resource "aws_api_gateway_resource" "taskapp_res_task" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_rest_api.taskapp.root_resource_id}"
  path_part   = "task"
}

module "task_list" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_task.id}"
  method      = "GET"
  path        = "${aws_api_gateway_resource.taskapp_res_task.path}"
  lambda      = "${var.apex_function_task_list}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
}

module "task_add" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_task.id}"
  method      = "POST"
  path        = "${aws_api_gateway_resource.taskapp_res_task.path}"
  lambda      = "${var.apex_function_task_add}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
}

resource "aws_api_gateway_deployment" "taskapp_deployment" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  stage_name  = "dev"
  description = "Deploy methods: ${module.task_list.http_method} ${module.task_add.http_method}"
}
