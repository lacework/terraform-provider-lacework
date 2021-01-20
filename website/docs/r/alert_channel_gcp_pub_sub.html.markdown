---
subcategory: "Alert Channels"
layout: "lacework"
page_title: "Lacework: lacework_alert_channel_gcp_pub_sub"
description: |-
  Create and manage GCP Pub Sub Alert Channel integrations
---

# lacework\_alert\_channel\_gcp\_pub\_sub

You can configure Lacework to forward events to this Google Cloud Pub/Sub asynchronous messaging service using the Lacework Google Cloud Pub/Sub alert channel.

To find more information see the [Lacework support documentation](https://support.lacework.com/hc/en-us/articles/360047496514-Google-Cloud-Pub-Sub).

## Example Usage

```hcl
resource "lacework_alert_channel_gcp_pub_sub" "example" {
  name       = "gcp-pub_sub"
  project_id = "lacework-191923"
  topic_id   = "lacework-alerts"
  credentials {
    client_id      = "123456789012345678900"
    client_email   = "email@abc-project-name.iam.gserviceaccount.com"
    private_key_id = "1234abcd1234abcd1234abcd1234abcd1234abcd"
    private_key    = "-----BEGIN PRIVATE KEY-----\n ... -----END PRIVATE KEY-----\n"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Alert Channel integration name.
* `project_id` - (Required) The name of the Gcp Project.
* `topic_id` - (Required) The id of the Gcp Pub Sub Topic.
* `credentials` - (Required) The credentials needed by the integration. See [Credentials](#credentials) below for details.
* `enabled` - (Optional) The state of the external integration. Defaults to `true`.

### Credentials

`credentials` supports the following arguments:

* `client_id` - (Required) The service account client ID.
* `client_email` - (Required) The service account client email.
* `private_key_id` - (Required) The service account private key ID.
* `private_key` - (Required) The service account private key.

## Import

A Lacework GCP Pub Sub Alert Channel integration can be imported using a `INT_GUID`, e.g.

```
$ terraform import lacework_alert_channel_gcp_pub_sub.data_export EXAMPLE_1234BAE1E42182964D23973F44CFEA3C4AB63B99E9A1EC5
```
-> **Note:** To retreive the `INT_GUID` from existing integrations in your account, use the
	Lacework CLI command `lacework integration list`. To install this tool follow
	[this documentation](https://github.com/lacework/go-sdk/wiki/CLI-Documentation#installation).
