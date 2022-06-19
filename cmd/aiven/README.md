aiven
=================

A tool to make ad-hoc [Aiven](https://aiven.io/) related Terraform state transformations.

Supports both Terraform and Terragrunt.

# Build & Run

## Build

1. Change directory to `cmd/aiven`
2. Run `go build`
3. Use created **aiven** binary for the next steps

## Run

Simply run the binary and use the provided prompt wizard to analyze existing configurations 
and make necessary migration transformations.

The tool assumes that you store your tfState.json file remotely on AWS.
That's why it checks beforehand validity of your AWS Credentials
and pulls then the state using correspondent terraform/terragrunt command.

If you only need to analyze your local tfState.json run `./aiven -state <path/to/tfState.json>`