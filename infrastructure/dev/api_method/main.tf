// Based on: https://github.com/TailorDev/hello-lambda/blob/master/api_method/main.tf

variable "rest_api_id" {
  description = "The ID of the associated REST API"
}

variable "resource_id" {
  description = "The API resource ID"
}

variable "method" {
  description = "The HTTP method"
  default     = "GET"
}

variable "path" {
  description = "The API resource path"
}

variable "lambda" {
  description = "The lambda name to invoke"
}

variable "region" {
  description = "The AWS region, e.g., eu-west-1"
}

variable "account_id" {
  description = "The AWS account ID"
}

variable "request_templates" {
  type = "map"
  description = "The request template map that maps query parameters to the lambda message"
  default = {}
}

variable "request_parameters" {
  // See: https://www.terraform.io/docs/providers/aws/r/api_gateway_method.html#request_parameters
  type = "map"
  description = "The request parameters for the method"
  default = {}
}

resource "aws_api_gateway_method" "request_method" {
  rest_api_id   = "${var.rest_api_id}"
  resource_id   = "${var.resource_id}"
  http_method   = "${var.method}"
  authorization = "NONE"

  request_parameters = "${var.request_parameters}"
}

resource "aws_api_gateway_integration" "request_method_integration" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_method.request_method.http_method}"
  type        = "AWS"
  uri         = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${var.lambda}/invocations"

  # AWS lambdas can only be invoked with the POST method
  integration_http_method = "POST"

  request_templates = "${var.request_templates}"
}

// Success

resource "aws_api_gateway_method_response" "response_method_200" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_integration.request_method_integration.http_method}"
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "response_method_integration_200" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_method_response.response_method_200.http_method}"
  status_code = "${aws_api_gateway_method_response.response_method_200.status_code}"
  response_templates = {
    "application/json" = ""
  }
}

// Bad Request

resource "aws_api_gateway_method_response" "response_method_400" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_integration_response.response_method_integration_200.http_method}"
  status_code = "400"

  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "response_method_integration_400" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_method_response.response_method_400.http_method}"
  status_code = "${aws_api_gateway_method_response.response_method_400.status_code}"
  selection_pattern = "BadRequest.*"
  response_templates = {
    "application/json" = ""
  }
}

// Internal Server Error

resource "aws_api_gateway_method_response" "response_method_500" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_integration_response.response_method_integration_400.http_method}"
  status_code = "500"

  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "response_method_integration_500" {
  rest_api_id = "${var.rest_api_id}"
  resource_id = "${var.resource_id}"
  http_method = "${aws_api_gateway_method_response.response_method_500.http_method}"
  status_code = "${aws_api_gateway_method_response.response_method_500.status_code}"
  selection_pattern = "InternalServerError.*"
  response_templates = {
    "application/json" = ""
  }
}

// Execution policy

resource "aws_lambda_permission" "allow_api_gateway" {
  function_name = "${var.lambda}"
  statement_id  = "AllowExecutionFromApiGateway"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.region}:${var.account_id}:${var.rest_api_id}/*/${var.method}${var.path}"
}

output "http_method" {
  value = "${aws_api_gateway_integration_response.response_method_integration_500.http_method}"
}