include "root" {
  path = find_in_parent_folders()
}

dependencies {
  paths = [
    "../PermissionSets/ExtendedReadOnlyAccess",
    "../PermissionSets/OrderEmployeeReportFilesAccess"
  ]
}

inputs = {
  name                  = "Order Processing Production"
  organizational_unit   = "Workloads_Prod"
  cost_center           = 1605
  komueb_product_ticket = "KOMUEB-1260"

  owner_email         = "firstname.lastname@domain.com"
  owner_jira_username = "firstname.lastname"

  group_permissions = {
    "pt-after_sales-oncall"           = ["AWSAdministratorAccess"]
    "pt-after_sales-order-processing" = ["ExtendedReadOnlyAccess"]
    "pt-order_employee_report-files"  = ["OrderEmployeeReportFilesAccess"]
  }
  personal_data_processed = true
  personal_data_stored    = true
}
