package reporting

import (
	"fmt"

	"github.com/pkg/errors"

	chef "github.com/chef/go-chef"
)

type CookbookVersion struct {
	Name    string
	Version string
}

func (cbv CookbookVersion) String() string {
	return fmt.Sprintf("%s(%s)", cbv.Name, cbv.Version)
}

type NodeReportItem struct {
	Name             string
	ChefVersion      string
	OS               string
	OSVersion        string
	CookbookVersions []CookbookVersion
}

type PartialSearchInterface interface {
	PartialExec(idx, statement string, params map[string]interface{}) (res chef.SearchResult, err error)
}

// TODO - should we offer this here as part of the struct, or allow the caller to handle it?
func (nri NodeReportItem) Array() []string {
	var cookbooks string
	for _, v := range nri.CookbookVersions {
		cookbooks = cookbooks + v.String() + " "
	}
	return []string{nri.Name,
		nri.ChefVersion,
		nri.OS,
		nri.OSVersion,
		cookbooks}
}

// func (nri NodeReportItem) String() string {
//   return ""
// }

// NOTE - we no longer need cfg. I'm not sure that this is best - I like having a single
//        cfg which includes the client, but did not want to create a full mock interface for
//        chef.client here - that belongs in go-chef, where it can be maintained alongside
//        any interface changes that originate there.
func Nodes(cfg *Reporting, searcher PartialSearchInterface) ([]NodeReportItem, error) {
	var (
		query = map[string]interface{}{
			"name":         []string{"name"},
			"chef_version": []string{"chef_packages", "chef", "version"},
			"os":           []string{"platform"},
			"os_version":   []string{"platform_version"},
			"cookbooks":    []string{"cookbooks"},
		}
	)

	pres, err := searcher.PartialExec("node", "*:*", query)
	if err != nil {
		e := errors.Wrap(err, "unable to get node(s) information")
		return nil, e
	}

	// We use len here and not pres.Rows, because when caller does not have permissions to
	// 	view all nodes in the result set, the actual returned number will be lower than
	// 	the value of Rows.

	results := make([]NodeReportItem, 0, len(pres.Rows))
	for _, element := range pres.Rows {

		// cookbook version arrives as [ NAME : { version: VERSION } - we extract that here.
		v := element.(map[string]interface{})["data"].(map[string]interface{})

		if v != nil {
			item := NodeReportItem{
				Name:        safeStringFromMap(v, "name"),
				OS:          safeStringFromMap(v, "os"),
				OSVersion:   safeStringFromMap(v, "os_version"),
				ChefVersion: safeStringFromMap(v, "chef_version"),
			}

			if v["cookbooks"] != nil {
				cookbooks := v["cookbooks"].(map[string]interface{})
				item.CookbookVersions = make([]CookbookVersion, 0, len(cookbooks))
				for k, v := range cookbooks {
					cbv := CookbookVersion{Name: k, Version: safeStringFromMap(v.(map[string]interface{}), "version")}
					item.CookbookVersions = append(item.CookbookVersions, cbv)
				}
			}
			results = append(results, item)
		}
	}
	return results, nil
}

// This returns the value referenced by `key` in `values`. If value is nil,
// it returns an empty string; otherwise it returns the original string.
// Will probably be more generally useful exposed as public in the reporting package.
func safeStringFromMap(values map[string]interface{}, key string) string {
	if values[key] == nil {
		return ""
	}
	return values[key].(string)
}
