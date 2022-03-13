module "komueb_3873_domain_control" {
  source = "../modules/aws-account"

  name                  = "Domain Control"
  organizational_unit   = local.organizational_unit_infrastructure_prod
  cost_center           = 1187
  komueb_product_ticket = "KOMUEB-3873"

  owner_email         = "heiko.rothe@idealo.de"
  owner_jira_username = "heiko.rothe"

  group_permissions = {
    "Cloud Shuttle" = ["AWSAdministratorAccess", "AWSReadOnlyAccess"]
  }
  user_permissions = {
    "heiko.rothe@idealo.de" = ["AWSPowerUserAccess"]
  }
}

module "komueb_3873_service_catalog_hub" {
  source = "../modules/aws-account"

  name                  = "Service Catalog Hub"
  organizational_unit   = local.organizational_unit_infrastructure_prod
  cost_center           = 1187
  komueb_product_ticket = "KOMUEB-3873"

  owner_email         = "moritz.brettschneider@idealo.de"
  owner_jira_username = "m.brettschneider"

  group_permissions = {
    "Cloud Shuttle" = ["AWSAdministratorAccess"]
  }
}

module "komueb_3873_aws_cost_management" {
  source = "../modules/aws-account"

  name                  = "AWS Cost Management"
  organizational_unit   = local.organizational_unit_workloads_prod
  cost_center           = 1187
  komueb_product_ticket = "KOMUEB-3873"

  owner_email         = "moritz.brettschneider@idealo.de"
  owner_jira_username = "m.brettschneider"

  group_permissions = {
    "Cloud Shuttle"   = [
      "AWSAdministratorAccess",
      module.controlling_access_permission_set.permission_set_name,
    ],
    "controlling-aws" = [
      module.controlling_access_permission_set.permission_set_name
    ],
    "pt-po-all" = [
      module.controlling_access_permission_set.permission_set_name
    ],
    "pt-headofs" = [
      module.controlling_access_permission_set.permission_set_name
    ],
    "pt-teamleads" = [
      module.controlling_access_permission_set.permission_set_name
    ]
  }
}

#All access for controlling
module "controlling_access_permission_set" {
  source                  = "../modules/aws-ssoadmin-permission-set"
  name                    = "ControllingCostReportAccess"
  ssoadmin_instance_arn   = local.ssoadmin_instance_arn
  inline_policy_document = data.aws_iam_policy_document.cost_controlling_access_policy.json
}

data "aws_iam_policy_document" "cost_controlling_access_policy" {
  statement {
    sid = "AthenaQueryExecAccess"
    actions   = [
      "athena:GetWorkGroup",
      "athena:GetQueryExecution",
      "athena:GetQueryResultsStream",
      "athena:GetQueryResults",
      "athena:ListQueryExecutions",
      "athena:ListNamedQueries",
      "athena:CreateNamedQuery",
      "athena:StartQueryExecution",
      "athena:StopQueryExecution"
    ]
    resources = [
      "arn:aws:athena:eu-central-1:173900957619:workgroup/primary"
    ]
  }
  statement {
    sid = "GlueReadAccess"
    actions   = [
      "glue:GetDatabase",
      "glue:GetDatabases",
      "glue:GetTable",
      "glue:GetTables",
      "glue:GetPartition",
      "glue:GetPartitions"
    ]
    resources = [
      "arn:aws:glue:eu-central-1:173900957619:catalog",
      "arn:aws:glue:eu-central-1:173900957619:database/athenacurcfn_general_cost_and_usage_report_idealo_prod",
      "arn:aws:glue:eu-central-1:173900957619:table/athenacurcfn_general_cost_and_usage_report_idealo_prod/cost_and_usage_data_status",
      "arn:aws:glue:eu-central-1:173900957619:table/athenacurcfn_general_cost_and_usage_report_idealo_prod/general_cost_and_usage_report_idealo_prod"
    ]
  }
  statement {
    sid = "S3DataSourceReadAccess"
    actions   = [
      "s3:ListBucket",
      "s3:GetObject"
    ]
    resources = [
      "arn:aws:s3:::957502001809-general-cost-and-usage-reports",
      "arn:aws:s3:::957502001809-general-cost-and-usage-reports/*"
    ]
  }
  statement {
    sid = "S3QueryResultsStorageAccess"
    actions   = [
      "s3:ListBucket",
      "s3:GetBucketLocation",
      "s3:GetObject",
      "s3:PutObject",
    ]
    resources = [
      "arn:aws:s3:::173900957619-athena-query-results",
      "arn:aws:s3:::173900957619-athena-query-results/*"
    ]
  }
  statement {
    sid = "S3CostReportsAccess"
    actions   = [
      "s3:ListBucket",
      "s3:GetObject"
    ]
    resources = [
      "arn:aws:s3:::173900957619-cost-reports-for-controlling/*",
      "arn:aws:s3:::173900957619-cost-reports-for-controlling",
      "arn:aws:s3:::173900957619-cost-reports-by-cost-center/*",
      "arn:aws:s3:::173900957619-cost-reports-by-cost-center"
    ]
  }
  statement {
    actions   = [
      "s3:ListAllMyBuckets"
    ]
    resources = [
      "*"
    ]
  }
}
