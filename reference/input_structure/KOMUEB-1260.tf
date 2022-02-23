module "komueb_1260_order_processing_production" {
  # reference to the module, does not need to be changed
  source = "../modules/aws-account"

  depends_on = [
    module.komueb_1260_order_employee_report_files_access_permission_set,
    module.komueb_1260_extended_read_only_access_permission_set
  ]

  name                  = "Order Processing Production"
  organizational_unit   = local.organizational_unit_workloads_prod
  cost_center           = 1605
  komueb_product_ticket = "KOMUEB-1260"

  owner_email         = "stefan.rudnitzki@idealo.de"
  owner_jira_username = "stefan.rudnitzki"

  personal_data_processed = true
  personal_data_stored    = true

  group_permissions = {
    "pt-after_sales-oncall"           = ["AWSAdministratorAccess", "AWSReadOnlyAccess"]
    "pt-after_sales-order-processing" = [
      module.komueb_1260_extended_read_only_access_permission_set.permission_set_name
    ]
    "pt-order_employee_report-files" = [
      module.komueb_1260_order_employee_report_files_access_permission_set.permission_set_name
    ]
  }
  user_permissions  = {
    "marcus.janke@idealo.de"    = [
      module.komueb_1260_extended_read_only_access_permission_set.permission_set_name
    ],
    "nicole.jaenchen@idealo.de" = [
      module.komueb_1260_extended_read_only_access_permission_set.permission_set_name
    ]
  }
}

module "komueb_1260_order_processing_staging" {
  # reference to the module, does not need to be changed
  source = "../modules/aws-account"

  name                  = "Order Processing Staging"
  organizational_unit   = local.organizational_unit_workloads_sdlc
  cost_center           = 1605
  komueb_product_ticket = "KOMUEB-1260"

  owner_email         = "stefan.rudnitzki@idealo.de"
  owner_jira_username = "stefan.rudnitzki"

  personal_data_processed = true
  personal_data_stored    = true

  group_permissions = {
    "pt-after_sales-oncall"           = ["AWSAdministratorAccess"]
    "pt-after_sales-order-processing" = ["AWSAdministratorAccess"]
  }
}

module "komueb_1260_extended_read_only_access_permission_set" {
  source                = "../modules/aws-ssoadmin-permission-set"
  name                  = "ExtendedReadOnlyAccess"
  ssoadmin_instance_arn = local.ssoadmin_instance_arn
  managed_policy_arns   = [
    "arn:aws:iam::aws:policy/ReadOnlyAccess"
  ]
  inline_policy_document = data.aws_iam_policy_document.extended_read_only_access.json
}

data "aws_iam_policy_document" "extended_read_only_access" {
  statement {
    actions = [
      "s3:GetObject",
      "sqs:ReceiveMessage",
      "lambda:GetFunction*",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
      "logs:GetLogRecord",
      "logs:StartQuery",
      "logs:StopQuery",
      "logs:GetQueryResults"
    ]
    resources = ["*"]
  }
}

module "komueb_1260_order_employee_report_files_access_permission_set" {
  source                = "../modules/aws-ssoadmin-permission-set"
  name                  = "OrderEmployeeReportFilesAccess"
  ssoadmin_instance_arn = local.ssoadmin_instance_arn
  managed_policy_arns   = [
    "arn:aws:iam::aws:policy/ReadOnlyAccess"
  ]
  inline_policy_document = data.aws_iam_policy_document.order_employee_report_bucket_access.json
}

data "aws_iam_policy_document" "order_employee_report_bucket_access" {
  statement {
    actions = [
      "s3:Get*",
      "s3:List*"
    ]
    resources = [
      "arn:aws:s3:::order-employee-report-files-production",
      "arn:aws:s3:::order-employee-report-files-production/*",
    ]
  }
}
