include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                   = "ControllingCostReportAccess"
  # language=JSON
  inline_policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
      "Sid": "AthenaQueryExecAccess",
      "Effect": "Allow",
      "Action": [
        "athena:GetWorkGroup",
        "athena:GetQueryExecution",
        "athena:GetQueryResultsStream",
        "athena:GetQueryResults",
        "athena:ListQueryExecutions",
        "athena:ListNamedQueries",
        "athena:CreateNamedQuery",
        "athena:StartQueryExecution",
        "athena:StopQueryExecution"
      ],
      "Resource": "arn:aws:athena:eu-central-1:173900957619:workgroup/primary"
    },
    {
      "Sid": "GlueReadAccess",
      "Effect": "Allow",
      "Action": [
        "glue:GetDatabase",
        "glue:GetDatabases",
        "glue:GetTable",
        "glue:GetTables",
        "glue:GetPartition",
        "glue:GetPartitions"
      ],
      "Resource": [
        "arn:aws:glue:eu-central-1:173900957619:catalog",
        "arn:aws:glue:eu-central-1:173900957619:database/athenacurcfn_general_cost_and_usage_report_idealo_prod",
        "arn:aws:glue:eu-central-1:173900957619:table/athenacurcfn_general_cost_and_usage_report_idealo_prod/cost_and_usage_data_status",
        "arn:aws:glue:eu-central-1:173900957619:table/athenacurcfn_general_cost_and_usage_report_idealo_prod/general_cost_and_usage_report_idealo_prod"
      ]
    },
    {
      "Sid": "S3DataSourceReadAccess",
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::957502001809-general-cost-and-usage-reports",
        "arn:aws:s3:::957502001809-general-cost-and-usage-reports/*"
      ]
    },
    {
      "Sid": "S3QueryResultsStorageAccess",
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetBucketLocation",
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": [
        "arn:aws:s3:::173900957619-athena-query-results",
        "arn:aws:s3:::173900957619-athena-query-results/*"
      ]
    },
    {
      "Sid": "S3CostReportsAccess",
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::173900957619-cost-reports-for-controlling/*",
        "arn:aws:s3:::173900957619-cost-reports-for-controlling",
        "arn:aws:s3:::173900957619-cost-reports-by-cost-center/*",
        "arn:aws:s3:::173900957619-cost-reports-by-cost-center"
      ]
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:ListAllMyBuckets",
      "Resource": "*"
    }
  ]
}
EOF
}
