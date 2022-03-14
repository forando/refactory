refactory
=================

A tool that automates ad-hoc Terraform project transformations.

Uses [hcl library](https://github.com/hashicorp/hcl) to parse
existing **.tf** files and generate new ones.


# Build & Run

## Local

1. Change directory to `cmd/refactory`
2. Run `go build`

# Open Questions

1. **depends_on** attribute in aws-account modules
2. **local.ssoadmin_instance_arn** in aws-ssoadmin-permission-set modules