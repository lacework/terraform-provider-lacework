//
// Author:: Darren Murray (<darren.murray@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// AlertRulesService is the service that interacts with
// the AlertRules schema from the Lacework APIv2 Server
type AlertRulesService struct {
	client *Client
}

// Valid inputs for AlertRule Source property
var AlertRuleSources = []string{"Agent", "AWS", "Azure", "GCP", "K8s", "OCI"}

// Valid inputs for AlertRule Categories property
var AlertRuleCategories = []string{"Anomaly", "Policy", "Composite"}

// Valid inputs for AlertRule SubCategories property
var AlertRuleSubCategories = []string{
	"Compliance",
	"Application",
	"Cloud Activity",
	"File",
	"Machine",
	"User",
	"Platform",
	"Kubernetes Activity",
	"Registry",
	"SystemCall",
	"Host Vulnerability",
	"Container Vulnerability",
	"Threat Intel",
	// Deprecated eventCategory values
	"App",
	"Cloud",
	"K8sActivity",
}

type alertRuleSeverity int

type AlertRuleSeverities []alertRuleSeverity

const AlertRuleEventType = "Event"

func (sevs AlertRuleSeverities) toInt() []int {
	var res []int
	for _, i := range sevs {
		res = append(res, int(i))
	}
	return res
}

func (sevs AlertRuleSeverities) ToStringSlice() []string {
	var res []string
	for _, i := range sevs {
		switch i {
		case AlertRuleSeverityCritical:
			res = append(res, "Critical")
		case AlertRuleSeverityHigh:
			res = append(res, "High")
		case AlertRuleSeverityMedium:
			res = append(res, "Medium")
		case AlertRuleSeverityLow:
			res = append(res, "Low")
		case AlertRuleSeverityInfo:
			res = append(res, "Info")
		default:
			continue
		}
	}
	return res
}

func NewAlertRuleSeverities(sevSlice []string) AlertRuleSeverities {
	var res AlertRuleSeverities
	for _, i := range sevSlice {
		sev := convertSeverity(i)
		if sev != AlertRuleSeverityUnknown {
			res = append(res, sev)
		}
	}
	return res
}

func NewAlertRuleSeveritiesFromIntSlice(sevSlice []int) AlertRuleSeverities {
	var res AlertRuleSeverities
	for _, i := range sevSlice {
		sev := convertSeverityInt(i)
		if sev != AlertRuleSeverityUnknown {
			res = append(res, sev)
		}
	}
	return res
}

func convertSeverity(sev string) alertRuleSeverity {
	switch strings.ToLower(sev) {
	case "critical":
		return AlertRuleSeverityCritical
	case "high":
		return AlertRuleSeverityHigh
	case "medium":
		return AlertRuleSeverityMedium
	case "low":
		return AlertRuleSeverityLow
	case "info":
		return AlertRuleSeverityInfo
	default:
		return AlertRuleSeverityUnknown
	}
}

func convertSeverityInt(sev int) alertRuleSeverity {
	switch sev {
	case 1:
		return AlertRuleSeverityCritical
	case 2:
		return AlertRuleSeverityHigh
	case 3:
		return AlertRuleSeverityMedium
	case 4:
		return AlertRuleSeverityLow
	case 5:
		return AlertRuleSeverityInfo
	default:
		return AlertRuleSeverityUnknown
	}
}

// Convert deprecated eventCatory values to subCategory values
func convertEventCategories(categories []string) []string {
	var res []string
	for _, c := range categories {
		switch c {
		case "App":
			res = append(res, "Application")
		case "Cloud":
			res = append(res, "Cloud Activity")
		case "K8sActivity":
			res = append(res, "Kubernetes Activity")
		default:
			res = append(res, c)
		}
	}
	return res
}

const (
	AlertRuleSeverityCritical alertRuleSeverity = 1
	AlertRuleSeverityHigh     alertRuleSeverity = 2
	AlertRuleSeverityMedium   alertRuleSeverity = 3
	AlertRuleSeverityLow      alertRuleSeverity = 4
	AlertRuleSeverityInfo     alertRuleSeverity = 5
	AlertRuleSeverityUnknown  alertRuleSeverity = 0
)

// NewAlertRule returns an instance of the AlertRule struct
//
// Basic usage: Initialize a new AlertRule struct, then
//
//	             use the new instance to do CRUD operations
//
//	  client, err := api.NewClient("account")
//	  if err != nil {
//	    return err
//	  }
//
//	  alertRule := api.NewAlertRule(
//			"Foo",
//			api.AlertRuleConfig{
//			Description: "My Alert Rule"
//			Severities: api.AlertRuleSeverities{api.AlertRuleSeverityHigh,
//			Channels: []string{"TECHALLY_000000000000AAAAAAAAAAAAAAAAAAAA"},
//			ResourceGroups: []string{"TECHALLY_111111111111AAAAAAAAAAAAAAAAAAAA"}
//	      },
//	    },
//	  )
//
//	  client.V2.AlertRules.Create(alertRule)
func NewAlertRule(name string, rule AlertRuleConfig) AlertRule {
	return AlertRule{
		Channels: rule.Channels,
		Type:     AlertRuleEventType,
		Filter: AlertRuleFilter{
			Name:               name,
			Enabled:            1,
			Description:        rule.Description,
			Severity:           rule.Severities.toInt(),
			ResourceGroups:     rule.ResourceGroups,
			AlertSubCategories: convertEventCategories(rule.AlertSubCategories),
			AlertCategories:    rule.AlertCategories,
			AlertSources:       rule.AlertSources,
		},
	}
}

func (rule AlertRuleFilter) Status() string {
	if rule.Enabled == 1 {
		return "Enabled"
	}
	return "Disabled"
}

// List returns a list of Alert Rules
func (svc *AlertRulesService) List() (response AlertRulesResponse, err error) {
	err = svc.client.RequestDecoder("GET", apiV2AlertRules, nil, &response)
	return
}

// Create creates a single Alert Rule
func (svc *AlertRulesService) Create(rule AlertRule) (
	response AlertRuleResponse,
	err error,
) {
	err = svc.client.RequestEncoderDecoder("POST", apiV2AlertRules, rule, &response)
	return
}

// Delete deletes a Alert Rule that matches the provided guid
func (svc *AlertRulesService) Delete(guid string) error {
	if guid == "" {
		return errors.New("specify an intgGuid")
	}

	return svc.client.RequestDecoder(
		"DELETE",
		fmt.Sprintf(apiV2AlertRuleFromGUID, guid),
		nil,
		nil,
	)
}

// Update updates a single Alert Rule of the provided guid.
func (svc *AlertRulesService) Update(data AlertRule) (
	response AlertRuleResponse,
	err error,
) {
	if data.Guid == "" {
		err = errors.New("specify a Guid")
		return
	}
	apiPath := fmt.Sprintf(apiV2AlertRuleFromGUID, data.Guid)
	err = svc.client.RequestEncoderDecoder("PATCH", apiPath, data, &response)
	return
}

// Get returns a raw response of the Alert Rule with the matching guid.
func (svc *AlertRulesService) Get(guid string, response interface{}) error {
	if guid == "" {
		return errors.New("specify a Guid")
	}
	apiPath := fmt.Sprintf(apiV2AlertRuleFromGUID, guid)
	return svc.client.RequestDecoder("GET", apiPath, nil, &response)
}

type AlertRuleConfig struct {
	Channels           []string
	Description        string
	Severities         AlertRuleSeverities
	ResourceGroups     []string
	AlertSubCategories []string
	AlertCategories    []string
	AlertSources       []string
}

type AlertRule struct {
	Guid     string          `json:"mcGuid,omitempty"`
	Type     string          `json:"type"`
	Channels []string        `json:"intgGuidList"`
	Filter   AlertRuleFilter `json:"filters"`
}

type AlertRuleFilter struct {
	Name                 string   `json:"name"`
	Enabled              int      `json:"enabled"`
	Description          string   `json:"description,omitempty"`
	Severity             []int    `json:"severity"`
	ResourceGroups       []string `json:"resourceGroups"`
	AlertSubCategories   []string `json:"subCategory"`
	AlertCategories      []string `json:"category"`
	AlertSources         []string `json:"source,omitempty"`
	CreatedOrUpdatedTime string   `json:"createdOrUpdatedTime,omitempty"`
	CreatedOrUpdatedBy   string   `json:"createdOrUpdatedBy,omitempty"`
}

type AlertRuleResponse struct {
	Data AlertRule `json:"data"`
}

type AlertRulesResponse struct {
	Data []AlertRule `json:"data"`
}
