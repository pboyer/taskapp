variable "account_id" {}
variable "access_key" {}
variable "secret_key" {}
variable "region" {
  default = "us-east-1"
}

variable "apex_function_task_add" {}
variable "apex_function_task_list" {}
variable "apex_function_task_update" {}
variable "apex_function_task_delete" {}
variable "apex_function_task_daily" {}

variable "apex_function_note_add" {}
variable "apex_function_note_get" {}
variable "apex_function_note_share" {}
variable "apex_function_note_unshare" {}
variable "apex_function_note_update" {}