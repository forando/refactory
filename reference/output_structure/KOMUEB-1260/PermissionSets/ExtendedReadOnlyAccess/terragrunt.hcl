include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                = "ExtendedReadOnlyAccess"
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
      "s3:GetObject",
      "sqs:ReceiveMessage",
      "lambda:GetFunction*",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
      "logs:GetLogRecord",
      "logs:StartQuery",
      "logs:StopQuery",
      "logs:GetQueryResults"
    ],
    "Resource": "*"
  }]
}
EOF
}
