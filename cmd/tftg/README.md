tftg (Terraform to Terragrunt)
=================

A tool that automates ad-hoc Terraform to Terragrunt project transformations.

Uses [hcl library](https://github.com/hashicorp/hcl) to parse
existing **.tf** files and generate new ones.


# Build & Run

## Build

1. Change directory to `cmd/tftg`
2. Run `go build`
3. Use created **tftg** binary for the next steps

## Run

1. Copy the **tftg** binary from Build stage into the working directory where all *.tf files located
2. Run `tftg [flags...] bootstrap` command. It will scaffold new Terragrunt structure
3. Set `export JIRA_USER=...` env variable
4. Set `export JIRA_PASSWORD=...` env variable
5. Configure proper AWS Credentials to use with aws cli pointing to the right account
6. Run `tftg [flags...] import` command that will fetch the current state of all necessary resources.

## Troubleshooting

### `to many open files` error
Run `ulimit -n` to see the limit.

Run `ulimit -n 4096` to make it larger.

### `terragrunt apply` on PermissionSets fails
Inline policy generation in json is currently malformed.
To fix it just reformat the file manually.

### Useful commands

If you need to delete some files of a particular pattern recursively from all directories
run `find . -name "*.bak" -type f -delete`