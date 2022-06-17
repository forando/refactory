include "root" {
  path = find_in_parent_folders()
}

inputs = {
  name                  = "Order Processing Staging"
  organizational_unit   = "Workloads_Sdlc"
  cost_center           = 1605
  komueb_product_ticket = "KOMUEB-1260"

  owner_email         = "firstname.lastname@domain.com"
  owner_jira_username = "firstname.lastname"

  group_permissions = {
    "pt-after_sales-oncall"           = ["AWSAdministratorAccess"]
    "pt-after_sales-order-processing" = ["AWSAdministratorAccess"]
  }
  personal_data_processed = true
  personal_data_stored    = true
}
