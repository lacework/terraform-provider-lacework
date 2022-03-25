Terraform Provider
==================

[![IaC](https://app.soluble.cloud/api/v1/public/badges/a1f83ada-bf2f-4029-a0bf-2cba01c2f548.svg)](https://app.soluble.cloud/repos/details/github.com/lacework/terraform-provider-lacework)  
[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/lacework/terraform-provider-lacework%2Ftest-build?type=cf-1&key=eyJhbGciOiJIUzI1NiJ9.NWVmNTAxOGU4Y2FjOGQzYTkxYjg3ZDEx.RJ3DEzWmBXrJX7m38iExJ_ntGv4_Ip8VTa-an8gBwBo)]( https://g.codefresh.io/pipelines/edit/new/builds?id=609b056f1e9a4249a520b52e&pipeline=test-build&projects=terraform-provider-lacework&projectId=609b049ae23d572127fccaff)

- Website: https://www.terraform.io
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin) 

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-lacework`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-lacework
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-lacework
$ make build
```

**Note**: For contributions created from forks, the repository should still be cloned under the `$GOPATH/src/github.com/terraform-providers/terraform-provider-lacework` directory to allow the provided `make` commands to properly run, build, and test this project.

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Further [usage documentation is available on the Terraform website](https://www.terraform.io/docs/providers/lacework/index.html).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.18+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-lacework
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
