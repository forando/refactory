include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                  = "Domain Control"
  organizational_unit   = "Infrastructure_Prod"
  cost_center           = 1187
  komueb_product_ticket = "KOMUEB-3873"

  owner_email         = "heiko.rothe@idealo.de"
  owner_jira_username = "heiko.rothe"

  group_permissions = {
    "Cloud Shuttle" = ["AWSAdministratorAccess", "AWSReadOnlyAccess"]
  }
}
