//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

import "time"

// SearchFilter is the representation of an advanced search payload
// for retrieving information out of the Lacework APIv2 Server
//
// An advanced example of a SearchFilter to search for an Agent
// Access Token that matches the provider token alias and return
// only the token found:
//
//		SearchFilter{
//			Filters: []Filter{
//				Filter{
//					Field:      "tokenAlias",
//					Expression: "eq",
//					Value:      "k8s-deployment,
//				},
//			},
//			Returns: []string{"accessToken"},
//		}
type SearchFilter struct {
	*TimeFilter `json:"timeFilter,omitempty"`
	Filters     []Filter `json:"filters,omitempty"`
	Returns     []string `json:"returns,omitempty"`
}

type Filter struct {
	Expression string   `json:"expression,omitempty"`
	Field      string   `json:"field,omitempty"`
	Value      string   `json:"value,omitempty"`
	Values     []string `json:"values,omitempty"`
}

type TimeFilter struct {
	StartTime *time.Time `json:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty"`
}
