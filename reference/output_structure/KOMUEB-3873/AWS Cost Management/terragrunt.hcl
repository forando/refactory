include "root" {
  path = find_in_parent_folders()
}

dependencies {
  paths = ["../PermissionSets/ControllingCostReportAccess"]
}

inputs = {
  name                  = "AWS Cost Management"
  organizational_unit   = "Workloads_Prod"
  cost_center           = 1187
  komueb_product_ticket = "KOMUEB-3873"

  owner_email         = "moritz.brettschneider@idealo.de"
  owner_jira_username = "m.brettschneider"

  group_permissions = {
    "pt-po-all"       = ["ControllingCostReportAccess"]
    "pt-headofs"      = ["ControllingCostReportAccess"]
    "pt-teamleads"    = ["ControllingCostReportAccess"]
    "Cloud Shuttle"   = [
      "AWSAdministratorAccess", "ControllingCostReportAccess"
    ]
    "controlling-aws" = ["ControllingCostReportAccess"]
  }
}
