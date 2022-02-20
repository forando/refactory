resource "aws_cloudformation_stack" "cur_athena_stack" {
  name = "CUR-AthenaStack"

  parameters = {
    SourceAccountID = local.org_management_account
    Stage = local.project_stage
  }

  capabilities = ["CAPABILITY_IAM"]

  template_body = file("${path.module}/stacks/cur-athena.yml")
}

resource "aws_s3_bucket" "athena_query_results" {
  bucket = "${data.aws_caller_identity.current.account_id}-athena-query-results"
  acl    = "private"

  versioning {
    enabled = true
  }

  logging {
    target_bucket = "${data.aws_caller_identity.current.account_id}-${data.aws_region.current.name}-s3-access-logs"
    target_prefix = "logs/"
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }

      bucket_key_enabled = true
    }
  }

  lifecycle_rule {
    id = "DeleteOldQueryResults"
    enabled = true

    expiration {
      days = 40
    }

    noncurrent_version_expiration {
      days = 3
    }
  }
}