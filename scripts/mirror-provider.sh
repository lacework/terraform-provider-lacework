#!/usr/bin/env bash

write_terraform_rc() {
  cat > ~/.terraformrc << EOL
provider_installation {
            filesystem_mirror {
              path    = "~/.terraform.d/plugins"
              include = ["lacework/lacework"]
            }
            direct {
              exclude = ["lacework/lacework"]
            }
          }
EOL
}

write_terraform_rc "$@" || exit 99
