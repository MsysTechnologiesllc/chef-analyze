//
// Copyright 2019 Chef Software, Inc.
// Author: Salim Afiune <afiune@chef.io>
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

package formatter

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chef/chef-analyze/pkg/reporting"
)

func MakeCookbooksReportCSV(state *reporting.CookbooksStatus) *FormattedResult {
	var (
		strBuilder strings.Builder
		errBuilder strings.Builder
		csvWriter  = csv.NewWriter(&strBuilder)
	)

	if state == nil || len(state.Records) == 0 {
		return &FormattedResult{"", ""}
	}

	tableHeaders := []string{"Cookbook Name", "Version"}
	if state.RunCookstyle {
		tableHeaders = append(tableHeaders,
			"File",
			"Offense",
			"Automatically Correctable",
			"Message",
		)
	}
	tableHeaders = append(tableHeaders, "Nodes")
	csvWriter.Write(tableHeaders)

	for _, record := range state.Records {
		nodesString := "None"
		if record.NumNodesAffected() != 0 {
			nodesString = strings.Join(record.Nodes, " ")
		}

		if state.RunCookstyle {
			for _, file := range record.Files {
				if len(file.Offenses) == 0 {
					continue
				}

				for _, offense := range file.Offenses {
					row := []string{
						record.Name,
						record.Version,
						file.Path,
						offense.CopName,
						"N",
						offense.Message,
						nodesString,
					}
					if offense.Correctable {
						row[4] = "Y"
					}
					csvWriter.Write(row)
				}
			}
		} else {
			row := []string{record.Name, record.Version, nodesString}
			csvWriter.Write(row)
		}

		for _, e := range record.Errors() {
			errBuilder.WriteString(fmt.Sprintf(" - %s (%s): %v\n", record.Name, record.Version, e))
		}
	}

	csvWriter.Flush()
	return &FormattedResult{strBuilder.String(), errBuilder.String()}
}

func MakeNodesReportCSV(records []*reporting.NodeReportItem) *FormattedResult {
	var (
		strBuilder strings.Builder
		errBuilder strings.Builder
		csvWriter  = csv.NewWriter(&strBuilder)
	)

	if len(records) == 0 {
		return &FormattedResult{"", ""}
	}

	tableHeaders := []string{"Node Name", "Chef Version", "Operating System", "Cookbooks"}
	csvWriter.Write(tableHeaders)

	for _, record := range records {
		var (
			cookbooksString = "None"
			cookbooksList   = record.CookbooksList()
		)
		if len(cookbooksList) != 0 {
			cookbooksString = strings.Join(cookbooksList, " ")
		}

		csvWriter.Write([]string{
			record.Name,
			record.ChefVersion,
			record.OSVersionPretty(),
			cookbooksString,
		})
	}

	csvWriter.Flush()
	return &FormattedResult{strBuilder.String(), errBuilder.String()}
}
