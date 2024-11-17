package manager

import (
	"sort"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/util"
)

func RegisterManagers() {}

func selectVersion(versions map[string][]core.Version, fileType string) (core.Version, bool) {
	var items []core.Version
	for _, vs := range versions {
		for _, v := range vs {
			if v.FileType == fileType {
				items = append(items, v)
			}
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		return util.CompareVersion(items[i].Version, items[j].Version) > 0
	})
	return util.Select(func(v core.Version) string { return v.Version }, items...)
}
