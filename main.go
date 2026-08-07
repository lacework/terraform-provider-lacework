package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/lacework/terraform-provider-lacework/lacework"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug: debug,
		// Required for -debug: this is the key Terraform matches against the
		// provider source address in TF_REATTACH_PROVIDERS. Without it the SDK
		// falls back to the literal "provider" and reattach never matches.
		ProviderAddr: "registry.terraform.io/lacework/lacework",
		ProviderFunc: lacework.Provider,
	})
}
