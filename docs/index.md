# AWS Xtras Provider

This provider aims to include resources and actions missing in the offcial AWS provider.

It will additionally create abstractions to commonly preformed actions. An example of this is vending accounts using AWS ControlTower.

This provider implements the official AWS provider's config, you would configure the provider exactly the same as you would config your AWS provider.

## Example Usage

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
