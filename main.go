package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/lacework/terraform-provider-lacework/lacework"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: lacework.Provider})
}
