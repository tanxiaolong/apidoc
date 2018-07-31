// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tanxiaolong/apidoc/types"
	"github.com/tanxiaolong/apidoc/vars"
)

type page struct {
	Title       string            `json:"title"`
	Version     string            `json:"version,omitempty"`
	BaseURL     string            `json:"baseURL"`
	LicenseName string            `json:"licenseName"`
	LicenseURL  string            `json:"licenseURL"`
	Content     string            `json:"content,omitempty"`
	Date        time.Time         `json:"date"`
	Elapsed     time.Duration     `json:"elapsed"`
	Groups      map[string]string `json:"groups"` // 组名与文件名的对应关系

	AppName    string `json:"appName"`
	AppURL     string `json:"appURL"`
	AppVersion string `json:"appVersion"`
}

type group struct {
	path string // 相对路径名

	Name string       `json:"name"` // 当前分组的名称
	Apis []*types.API `json:"apis"` // 当前分组的 api 文档
}

func render(docs *types.Doc, opt *Options) error {
	groups := make(map[string]*group, 100)

	for _, api := range docs.Apis {
		name := strings.ToLower(api.Group)
		path := filepath.Join(opt.dataDir, vars.GroupFilePrefix+name+".json")

		if groups[name] == nil {
			groups[name] = &group{
				path: path,
				Name: api.Group, // 名称区分大小写，不采用 name 变量
				Apis: make([]*types.API, 0, 100),
			}
		}
		groups[name].Apis = append(groups[name].Apis, api)
	}

	names := make(map[string]string, len(groups))
	for name, group := range groups {
		names[group.Name] = path.Join(vars.JSONDataDirName, vars.GroupFilePrefix+name+".json")

		// 排序
		sort.SliceStable(group.Apis, func(i, j int) bool {
			return group.Apis[i].URL < group.Apis[j].URL
		})
	}

	page := &page{
		Title:       docs.Title,
		Version:     docs.Version,
		BaseURL:     docs.BaseURL,
		LicenseName: docs.LicenseName,
		LicenseURL:  docs.LicenseURL,
		Content:     docs.Content,
		Date:        time.Now(),
		Elapsed:     opt.Elapsed,
		Groups:      names,

		AppName:    vars.Name,
		AppURL:     vars.OfficialURL,
		AppVersion: vars.Version(),
	}

	if err := renderPage(page, opt.dataDir); err != nil {
		return err
	}

	return renderGroups(groups, opt)
}

func renderPage(p *page, destDir string) error {
	path := filepath.Join(destDir, vars.PageFileName+".json")

	return renderJSON(p, path)
}

func renderGroups(groups map[string]*group, o *Options) error {
	for _, g := range groups {
		if !o.groupIsEnable(g.Name) {
			continue
		}

		if err := renderJSON(g, g.path); err != nil {
			return err
		}
	}

	return nil
}

func renderJSON(obj interface{}, path string) error {
	data, err := json.MarshalIndent(obj, "", strings.Repeat(" ", vars.JSONIndent))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}
