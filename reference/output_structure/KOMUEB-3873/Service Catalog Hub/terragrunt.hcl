include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                  = "Service Catalog Hub"
  organizational_unit   = "Infrastructure_Prod"
  cost_center           = 1187
  komueb_product_ticket = "KOMUEB-3873"

  owner_email         = "moritz.brettschneider@idealo.de"
  owner_jira_username = "m.brettschneider"

  group_permissions = {
    "Cloud Shuttle" = ["AWSAdministratorAccess"]
  }
}
