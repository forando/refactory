include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                = "OrderEmployeeReportFilesAccess"
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/ReadOnlyAccess"
  ]
  # language=JSON
  inline_policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "",
    "Effect": "Allow",
    "Action": [
      "s3:Get*",
      "s3:List*"
    ],
    "Resource": [
      "arn:aws:s3:::order-employee-report-files-production",
      "arn:aws:s3:::order-employee-report-files-production/*"
    ]
  }]
}
EOF
}
