<img src="https://techally-content.s3-us-west-1.amazonaws.com/public-content/lacework_logo_full.png" width="600">

Terraform Provider for Lacework
==================

[![IaC](https://app.soluble.cloud/api/v1/public/badges/a1f83ada-bf2f-4029-a0bf-2cba01c2f548.svg)](https://app.soluble.cloud/repos/details/github.com/lacework/terraform-provider-lacework)  
[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/lacework/terraform-provider-lacework%2Ftest-build?type=cf-1&key=eyJhbGciOiJIUzI1NiJ9.NWVmNTAxOGU4Y2FjOGQzYTkxYjg3ZDEx.RJ3DEzWmBXrJX7m38iExJ_ntGv4_Ip8VTa-an8gBwBo)]( https://g.codefresh.io/pipelines/edit/new/builds?id=609b056f1e9a4249a520b52e&pipeline=test-build&projects=terraform-provider-lacework&projectId=609b049ae23d572127fccaff)

- Website: https://www.terraform.io
- Registry: https://registry.terraform.io/providers/lacework/lacework/latest

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.15.x
-	[Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/lacework/terraform-provider-lacework`

```sh
mkdir -p $GOPATH/src/github.com/lacework; cd $GOPATH/src/github.com/lacework
git clone git@github.com:lacework/terraform-provider-lacework
```

Enter the provider directory, prepare and build the provider

```sh
cd $GOPATH/src/github.com/lacework/terraform-provider-lacework
make prepare
make build
```

**Note**: For contributions created from forks, the repository should still be cloned under the `$GOPATH/src/github.com/lacework/terraform-provider-lacework` directory to allow the provided `make` commands to properly run, build, and test this project.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version
1.18+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well
as adding `$GOPATH/bin` to your `$PATH`.

To install the provider, run `make install`. This will build the provider, create the `$HOME/.terraformrc` and put the
provider binary in the `$HOME/.terraform.d/plugins` directory.

```sh
make install
```

From here, you can `cd` into any folder that contains Terraform code and run `terraform init`, use one of the examples
in the [examples/](examples/) folder.

```sh
cd examples/data_source_lacework_api_token
terraform init
terraform apply
```

In order run the providers' unit tests, run the command `make test`:

```sh
make test
```

In order to run the full suite of integration tests, run `make integration-test`.

*Note:* Integration tests create real resources so you need to have a Lacework environment.

```sh
LW_ACCOUNT="<YOUR_ACCOUNT>" \
  LW_API_KEY="<YOUR_API_KEY>" \
  LW_API_SECRET="<YOUR_API_SECRET>" \
  make integration-test
```
### Running Specific Integration Tests (RegEx)
When working on new tests or existing tests, you can use a regex to run only specific integration tests. For example,
to run only the tests related to the resource `lacework_policy` use the regex `TestPolicy`:

```sh
make integration-test regex=TestPolicy
```

Uninstall Developer Environment
---------------------------

If you are doing development of the Lacework provider and you try to use the provider we release to the Terraform
Registry, you won't be able to use it and you will see an error message similar to:

```
│ Error: Failed to query available provider packages
│
│ Could not retrieve the list of available versions for provider lacework/lacework: no available releases match the given constraints
```

To fix this issue you need to uninstall the local provider (your developer environment) with the command:

```sh
make uninstall
```
