provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

# DynamoDB

resource "aws_dynamodb_table" "taskapp_tasks" {
  name           = "taskapp_tasks"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "taskapp_notes" {
  name           = "taskapp_notes"
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
  description = "A simple task list app"
}

resource "aws_api_gateway_resource" "taskapp_res_task" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_rest_api.taskapp.root_resource_id}"
  path_part   = "task"
}

resource "aws_api_gateway_resource" "taskapp_res_task_list" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_task.id}"
  path_part   = "list"
}

// GET /task/list?user=foo&priority=12&description=foo&complete=date
module "task_list" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_task_list.id}"
  method      = "GET"
  path        = "${aws_api_gateway_resource.taskapp_res_task_list.path}"
  lambda      = "${var.apex_function_task_list}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #foreach($param in $input.params().querystring.keySet())
    #if($param == "priority")
      "$param": $util.escapeJavaScript($input.params().querystring.get($param))
    #{else}
      "$param": "$util.escapeJavaScript($input.params().querystring.get($param))" 
    #end
    #if($foreach.hasNext),#end
  #end
}
EOF
  }
}

// POST /task
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

// PUT /task/{updateTaskId}
resource "aws_api_gateway_resource" "taskapp_res_task_update" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_task.id}"
  path_part   = "{updateTaskId}"
}

module "task_update" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_task_update.id}"
  method      = "PUT"
  path        = "${aws_api_gateway_resource.taskapp_res_task_update.path}"
  lambda      = "${var.apex_function_task_update}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #if($input.params('updateTaskId'))"id" : "$input.params('updateTaskId')"#end
}
EOF
  }
}

// DELETE /task/{deleteTaskId}
resource "aws_api_gateway_resource" "taskapp_res_task_delete" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_task.id}"
  path_part   = "{deleteTaskId}"
}

module "task_delete" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_task_delete.id}"
  method      = "DELETE"
  path        = "${aws_api_gateway_resource.taskapp_res_task_delete.path}"
  lambda      = "${var.apex_function_task_delete}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #if($input.params('deleteTaskId'))"id" : "$input.params('deleteTaskId')"#end
}
EOF
  }
}

resource "aws_api_gateway_deployment" "taskapp_deployment" {
  depends_on = ["module.task_delete", "module.task_update", "module.task_add", "module.task_list"]

  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  stage_name  = "dev"
  description = "Deploy methods"
}
