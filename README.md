# Duvet

Simple tool to compare available AWS API create methods against their usage in the Terraform AWS Provider. Generally 
Terraform resources correspond to a CreateXXX API within the AWS Go SDK though this is not always the case. 

The tool does this by:
- Creating clients for all available AWS Go SDK service packages. 
- Reflect over the available methods and keep track of those prefixed with Create.
- Walk `*.go` files in the corresponding `internal/service/${serviceName}` directory to look for usage of those Create calls 
- Display the data as html or markdown.

## Requirements

- Go
- Locally checked out copies of both the aws-go-sdk (v2) and the Terraform AWS Provider.