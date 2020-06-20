# controltower_account_vending Resource

Vend an ControlTower account. This resources abstracts some of the Service Catalog complexities.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the account to be vended. This will also be used for the provisione product in Service Catalog
* `product_id` - (Required) The product Id for the ControlTower account vending product.
* `artefact_id` - (Required) The artefact Id to provision (version of the product)
* `parameters` - (Required) List attributes needed to provision the product. See the example of the parameters needed for a ControlTower account.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Arn of the provisioned product
* `created_time` - Created Time
* `record_id` - The record Id of the product provisioning
* `account_id` - The account Id of the newly provisioned account.
