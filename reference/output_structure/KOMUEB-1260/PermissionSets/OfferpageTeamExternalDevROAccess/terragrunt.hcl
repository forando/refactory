include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                   = "OfferpageTeamExternalDevROAccess"
  # language=JSON
  inline_policy_document = <<EOF
{
  "Statement": {
    "Action": ["secretsmanager:GetResourcePolicy"]
  }
}
EOF
}
