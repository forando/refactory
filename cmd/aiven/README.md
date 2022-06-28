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

The tool assumes that you store your terraform.tfstate file remotely on AWS.
That's why it checks beforehand validity of your AWS Credentials
and pulls then the state using correspondent terraform/terragrunt command.

This tool helps you to:
- remove PeeringConnections from your state;
- move ConnectionAccepter resources under the new module control

If you need to analyze your local tfState.json files run `./aiven <path/to/terraform1.tfstate> <path/to/terraform2.tfstate>`.
This way you'll be prompted to import resources from those files into your current state.

### Terraform
- Run `terraform state pull > path/to/terraform.tfstate` for having the old state locally as a backup
- Run `terraform plan` and ensure that the plan is empty
- Run `./aiven` and simply follow the prompt instructions
- Rerun `terraform plan` and ensure that the plan is still empty

### Terragrunt
The problem is that terragrunt splits the project into multiple states (one terraform.tfstate file for each new module).
That means that instead of moving the resources within *the same state* you would have to drop them in the current state and *reimport* them *into the new module state*.

- Run `terragrunt state pull > path/to/terraform.tfstate` from within an old terragrunt configuration folder
- Run `./aiven path/to/terraform.tfstate` and follow the prompt guide all the way down to the new module generation and hit exit afterwards
- This new module is in terraform format so, using the data from it, you would need to make a valid terragrunt module and configuration
- Rerun `./aiven path/to/terraform.tfstate` (do not generate a new module again) and proceed with importing resources into the new state
- Run `terragrunt plan` from within the new terragrunt configuration folder and ensure that the plan is empty
- Get rid of the old terragrunt configuration folder and corresponding **terraform.tfstate** file in an AWS **S3 bucket**