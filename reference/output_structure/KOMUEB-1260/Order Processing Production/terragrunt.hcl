include "root" {
  path = find_in_parent_folders()
}

dependencies {
  paths = ["../PermissionSets/OfferpageTeamExternalDevROAccess"]
}

inputs = {
  name                  = "Order Processing Production"
  organizational_unit   = "Workloads_Prod"
  cost_center           = 1605
  komueb_product_ticket = "KOMUEB-1260"

  owner_email         = "stefan.rudnitzki@idealo.de"
  owner_jira_username = "stefan.rudnitzki"

  group_permissions = {
    "pt-after_sales-oncall" = ["AWSAdministratorAccess", "AWSReadOnlyAccess", "OfferpageTeamExternalDevROAccess"]
  }

  user_permissions = {
    "marcus.janke@idealo.de" = ["ExtendedReadOnlyAccess"]
  }
}
