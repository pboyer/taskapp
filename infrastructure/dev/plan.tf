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

# Task Daily event

resource "aws_cloudwatch_event_rule" "daily" {
    name = "daily"
    description = "Fires every five minutes"
    schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "task_daily" {
    rule = "${aws_cloudwatch_event_rule.daily.name}"
    arn = "${var.apex_function_task_daily}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_task" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${var.apex_function_task_daily}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.daily.arn}"
}

# API Gateway

resource "aws_api_gateway_rest_api" "taskapp" {
  name = "Task App"
  description = "A simple task list app"
}

# Task resources ---------------------------------------------------------------

resource "aws_api_gateway_resource" "taskapp_res_task" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_rest_api.taskapp.root_resource_id}"
  path_part   = "task"
}

# GET /task/list?user=foo&priority=12&description=foo&complete=date

resource "aws_api_gateway_resource" "taskapp_res_task_list" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_task.id}"
  path_part   = "list"
}

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

# POST /task

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

# PUT /task/{updateTaskId}

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

# DELETE /task/{deleteTaskId}

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

# Note resources ---------------------------------------------------------------

resource "aws_api_gateway_resource" "taskapp_res_note" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_rest_api.taskapp.root_resource_id}"
  path_part   = "note"
}

# POST /note

module "note_add" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_note.id}"
  method      = "POST"
  path        = "${aws_api_gateway_resource.taskapp_res_note.path}"
  lambda      = "${var.apex_function_note_add}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
}

# GET /note/{getNoteId}

resource "aws_api_gateway_resource" "taskapp_res_note_get" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_note.id}"
  path_part   = "{getNoteId}"
}

module "note_get" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_note_get.id}"
  method      = "GET"
  path        = "${aws_api_gateway_resource.taskapp_res_note_get.path}"
  lambda      = "${var.apex_function_note_get}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #if($input.params('getNoteId'))"id" : "$input.params('getNoteId')"#end
}
EOF
  }
}

# PUT /note/{updateNoteId}

resource "aws_api_gateway_resource" "taskapp_res_note_update" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_note.id}"
  path_part   = "{updateNoteId}"
}

module "note_update" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_note_update.id}"
  method      = "PUT"
  path        = "${aws_api_gateway_resource.taskapp_res_note_update.path}"
  lambda      = "${var.apex_function_note_update}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #if($input.params('updateNoteId'))"id" : "$input.params('updateNoteId')"#end
}
EOF
  }
}

# PUT /note/share/{shareNoteId}

resource "aws_api_gateway_resource" "taskapp_res_note_share" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_note.id}"
  path_part   = "share"
}

resource "aws_api_gateway_resource" "taskapp_res_note_share_id" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_note_share.id}"
  path_part   = "{shareNoteId}"
}

module "note_share" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_note_share_id.id}"
  method      = "PUT"
  path        = "${aws_api_gateway_resource.taskapp_res_note_share_id.path}"
  lambda      = "${var.apex_function_note_share}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #if($input.params('shareNoteId'))"id" : "$input.params('shareNoteId')"#end
}
EOF
  }
}

# PUT /note/unshare/{unshareNoteId}

resource "aws_api_gateway_resource" "taskapp_res_note_unshare" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_note.id}"
  path_part   = "unshare"
}

resource "aws_api_gateway_resource" "taskapp_res_note_unshare_id" {
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  parent_id   = "${aws_api_gateway_resource.taskapp_res_note_unshare.id}"
  path_part   = "{unshareNoteId}"
}

module "note_unshare" {
  source      = "./api_method"
  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  resource_id = "${aws_api_gateway_resource.taskapp_res_note_unshare_id.id}"
  method      = "PUT"
  path        = "${aws_api_gateway_resource.taskapp_res_note_unshare_id.path}"
  lambda      = "${var.apex_function_note_unshare}"
  region      = "${var.region}"
  account_id  = "${var.account_id}"
  request_templates = {
  "application/json" = <<EOF
{
  #if($input.params('unshareNoteId'))"id" : "$input.params('unshareNoteId')"#end
}
EOF
  }
}

# Stage deploy

resource "aws_api_gateway_deployment" "taskapp_deployment" {
  depends_on = ["module.note_add", "module.note_get", "module.note_update", "module.note_share", "module.note_unshare", 
                "module.note_unshare", "module.task_delete", "module.task_update", "module.task_add", "module.task_list"]

  rest_api_id = "${aws_api_gateway_rest_api.taskapp.id}"
  stage_name  = "dev"
  description = "Deploy methods"
}
