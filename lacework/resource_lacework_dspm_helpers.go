package lacework

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/v2/api"
)

// buildDspmProps extracts DSPM config fields from the Terraform schema and
// builds an api.DspmProps. Returns nil if none of the optional fields are set.
func buildDspmProps(d *schema.ResourceData) (*api.DspmProps, error) {
	var cfg api.DspmPropsConfig
	hasProps := false

	if v, ok := d.GetOk("scan_frequency_hours"); ok {
		hours := v.(int)
		cfg.ScanIntervalHours = &hours
		hasProps = true
	}

	if v, ok := d.GetOk("max_file_size_mb"); ok {
		mb := v.(int)
		bytes := mb * 1024 * 1024
		cfg.MaxDownloadBytes = &bytes
		hasProps = true
	}

	if v, ok := d.GetOk("datastore_filters"); ok {
		filters := v.([]interface{})
		if len(filters) > 0 && filters[0] != nil {
			filterMap := filters[0].(map[string]interface{})
			filterMode := filterMap["filter_mode"].(string)
			dspmFilter := &api.DspmDatastoreFilters{
				FilterMode: filterMode,
			}

			names := castAndTransformStringSlice(filterMap["datastore_names"].([]interface{}), func(s string) string { return s })
			if filterMode == "ALL" {
				if len(names) > 0 {
					return nil, fmt.Errorf("datastore_names must not be set when filter_mode is 'ALL'")
				}
			} else {
				if len(names) == 0 {
					return nil, fmt.Errorf("datastore_names is required when filter_mode is '%s'", filterMode)
				}
				dspmFilter.DatastoreNames = names
			}

			cfg.DatastoreFilters = dspmFilter
			hasProps = true
		}
	}

	if !hasProps {
		return nil, nil
	}

	return &api.DspmProps{
		Dspm: &cfg,
	}, nil
}

// updateDspmStatus calls the DSPM status API to mark the integration as
// configured when datastore_filters are present. When filters are absent the
// UI shows "Setup Required" based on the missing DATASTORE_FILTERS prop alone,
// so no status call is needed.
func updateDspmStatus(d *schema.ResourceData, client *api.Client, serverToken string) error {
	if _, ok := d.GetOk("datastore_filters"); ok {
		return client.V2.CloudAccounts.UpdateDspmStatus(serverToken, api.DspmStatusRequest{
			Ok:      true,
			Message: "SUCCESS",
		})
	}
	return nil
}

// readDspmProps reads DSPM props from an API response and sets them on the
// Terraform schema.
func readDspmProps(d *schema.ResourceData, props *api.DspmProps) {
	if props == nil || props.Dspm == nil {
		return
	}

	cfg := props.Dspm

	if cfg.ScanIntervalHours != nil {
		d.Set("scan_frequency_hours", *cfg.ScanIntervalHours)
	}

	if cfg.MaxDownloadBytes != nil {
		d.Set("max_file_size_mb", *cfg.MaxDownloadBytes/(1024*1024))
	}

	if cfg.DatastoreFilters != nil {
		filter := map[string]interface{}{
			"filter_mode": cfg.DatastoreFilters.FilterMode,
		}
		if cfg.DatastoreFilters.DatastoreNames != nil {
			filter["datastore_names"] = cfg.DatastoreFilters.DatastoreNames
		}
		d.Set("datastore_filters", []map[string]interface{}{filter})
	}
}
