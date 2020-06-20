Terraform Provider Scaffolding
==================

This repository is a *template* for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

 - A resource, and a data source (`internal/provider/`),
 - Documentation (`website/`),
 - Miscellanious meta files.
 
These files contain boilerplate code that you will need to edit to create your own Terraform provider. A full guide to creating Terraform providers can be found at [Writing Custom Providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.12

Building The Provider
---------------------

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

Adding Dependencies
---------------------

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.


Using the provider
----------------------

```hcl
provider awsx {
  region = "eu-west-1"
}

resource controltower_account_vending crypto {

  product_id  = "prod-nl7pbqs2n3rjy"
  artefact_id = "pa-htxzmae7h7bd2"

  parameters = {
    SSOUserFirstName          = "RootName"
    SSOUserLastName           = "RootSurname"
    SSOUserEmail              = "me+awsx-provider-1-SSO@gmail.com"
    AccountEmail              = "me+awsx-provider-1-Account@gmail.com"
    ManagedOrganizationalUnit = "Custom"
    AccountName               = "awsx-provider-1"

  }

  name = "controltower-provider-1"

  # This will prevent things going sideways if you dynamically 
  # lookup the latest artefact id
  lifecycle {
    ignore_changes = [
      artefact_id,
    ]
  }
}

output account_id {
  value = controltower_account_vending.crypto.account_id
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

Publishing Providers to the Terraform Registry
---------------------------
https://docs.google.com/document/d/1J4p5KFH129wZbSF0XUioDbHSB7uZA-l2o1psLa-3YS4/edit#