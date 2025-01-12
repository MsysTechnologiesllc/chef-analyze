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

package formatter_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	subject "github.com/chef/chef-analyze/pkg/formatter"
	"github.com/chef/chef-analyze/pkg/reporting"
)

func TestMakeCookbooksReportCSV_Nil(t *testing.T) {
	assert.Equal(t,
		&subject.FormattedResult{Report: "", Errors: ""},
		subject.MakeCookbooksReportCSV(nil))
}

func TestMakeCookbooksReportCSV_NoRecords(t *testing.T) {
	var expected subject.FormattedResult // empty result
	var cbStatus reporting.CookbooksStatus
	actual := subject.MakeCookbooksReportCSV(&cbStatus)
	assert.Equal(t, expected, *actual)
}

func TestMakeCookbooksReportCSV_WithUnverifiedRecords(t *testing.T) {
	cbStatus := reporting.CookbooksStatus{
		RunCookstyle: false,
		Records: []*reporting.CookbookRecord{
			&reporting.CookbookRecord{Name: "my-cookbook", Version: "1.0", Nodes: []string{"node-1", "node-2"}},
		},
	}

	lines := strings.Split(subject.MakeCookbooksReportCSV(&cbStatus).Report, "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "Cookbook Name,Version,Nodes", lines[0])
	assert.Equal(t, "my-cookbook,1.0,node-1 node-2", lines[1])
	assert.Equal(t, "", lines[2])
}

func TestMakeCookbooksReportCSV_WithVerifiedRecords(t *testing.T) {
	cbStatus := reporting.CookbooksStatus{
		RunCookstyle: true,
		Records: []*reporting.CookbookRecord{
			&reporting.CookbookRecord{Name: "my-cookbook", Version: "1.0", Nodes: []string{"node-1", "node-2"},
				Files: []reporting.CookbookFile{
					reporting.CookbookFile{Path: "/path/to/file.rb",
						Offenses: []reporting.CookstyleOffense{
							reporting.CookstyleOffense{CopName: "ChefDeprecations/Blah", Message: "some description", Correctable: true},
						}}}}}}

	actual := subject.MakeCookbooksReportCSV(&cbStatus)
	lines := strings.Split(actual.Report, "\n")
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "Cookbook Name,Version,File,Offense,Automatically Correctable,Message,Nodes", lines[0])
	assert.Equal(t, "my-cookbook,1.0,/path/to/file.rb,ChefDeprecations/Blah,Y,some description,node-1 node-2", lines[1])
	assert.Equal(t, "", lines[2])
}

func TestMakeCookbooksReportCSV_ErrorReport(t *testing.T) {
	cbStatus := reporting.CookbooksStatus{
		Records: []*reporting.CookbookRecord{
			&reporting.CookbookRecord{Name: "my-cookbook", Version: "1.0", DownloadError: errors.New("could not download")},
			&reporting.CookbookRecord{Name: "their-cookbook", Version: "1.1", UsageLookupError: errors.New("could not look up usage")},
			&reporting.CookbookRecord{Name: "our-cookbook", Version: "1.2", CookstyleError: errors.New("cookstyle error")},
		},
	}

	actual := subject.MakeCookbooksReportCSV(&cbStatus)
	lines := strings.Split(actual.Errors, "\n")
	assert.Equal(t, 4, len(lines))
	assert.Equal(t, lines[0], " - my-cookbook (1.0): could not download")
	assert.Equal(t, lines[1], " - their-cookbook (1.1): could not look up usage")
	assert.Equal(t, lines[2], " - our-cookbook (1.2): cookstyle error")
	assert.Equal(t, lines[3], "")
}

func TestMakeNodesReportCSV_Nil(t *testing.T) {
	assert.Equal(t,
		&subject.FormattedResult{Report: "", Errors: ""},
		subject.MakeNodesReportCSV(nil))
}

func TestMakeNodesReportCSV_NoRecords(t *testing.T) {
	var expected subject.FormattedResult // empty result
	var nodesReport = []*reporting.NodeReportItem{}
	actual := subject.MakeNodesReportCSV(nodesReport)
	assert.Equal(t, expected, *actual)
}

func TestMakeNodesReportCSV_WithRecords(t *testing.T) {
	nodesReport := []*reporting.NodeReportItem{
		&reporting.NodeReportItem{Name: "node1", ChefVersion: "12.22", OS: "windows", OSVersion: "10.1",
			CookbookVersions: []reporting.CookbookVersion{
				reporting.CookbookVersion{Name: "mycookbook", Version: "1.0"}},
		},
		&reporting.NodeReportItem{Name: "node2", ChefVersion: "13.11", OS: "", OSVersion: "",
			CookbookVersions: []reporting.CookbookVersion{
				reporting.CookbookVersion{Name: "mycookbook", Version: "1.0"},
				reporting.CookbookVersion{Name: "test", Version: "9.9"},
			},
		},
		&reporting.NodeReportItem{Name: "node3", ChefVersion: "15.00", OS: "ubuntu", OSVersion: "16.04",
			CookbookVersions: nil},
	}

	lines := strings.Split(subject.MakeNodesReportCSV(nodesReport).Report, "\n")
	if assert.Equal(t, 5, len(lines)) {
		assert.Equal(t, "Node Name,Chef Version,Operating System,Cookbooks", lines[0])
		assert.Equal(t, "node1,12.22,windows v10.1,mycookbook(1.0)", lines[1])
		assert.Equal(t, "node2,13.11,,mycookbook(1.0) test(9.9)", lines[2])
		assert.Equal(t, "node3,15.00,ubuntu v16.04,None", lines[3])
		assert.Equal(t, "", lines[4])
	}
}
