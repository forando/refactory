remote_state {
  backend  = "s3"
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite"
  }
  config = {
    bucket         = "idealo-test-org-tg-state"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = "eu-central-1"
    encrypt        = true
    dynamodb_table = "idealo-test-org-tg-lock"
  }
}

terraform {
  source = length(regexall(".*\\/PermissionSets\\/.*", get_original_terragrunt_dir())) > 0 ? "${get_path_to_repo_root()}/modules//aws-ssoadmin-permission-set" : "${get_path_to_repo_root()}/modules//aws-account"

  before_hook "link_lock_file" {
    commands = ["init"]
    execute  = [
      "ln", "-sf",
      length(regexall(".*\\/PermissionSets\\/.*", get_original_terragrunt_dir())) > 0 ? "${get_path_to_repo_root()}/modules/aws-ssoadmin-permission-set/.terraform.lock.hcl" : "${get_path_to_repo_root()}/modules/aws-account/.terraform.lock.hcl",
      "${get_original_terragrunt_dir()}/.terraform.lock.hcl"
    ]
  }
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"

  # language=HCL
  contents = <<EOF
provider "aws" {
  region = "eu-central-1"
  allowed_account_ids = ["573275350257"]
}

%{ if length(regexall(".*\\/PermissionSets\\/.*", get_original_terragrunt_dir())) == 0 ~}
provider "controltower" {
  region = "eu-central-1"
}

provider "jira" {
  url = "https://jira.eu.idealo.com/issues"
}
%{ endif ~}
EOF
}
