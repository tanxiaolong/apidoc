// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"os"
	"path/filepath"
	"time"

	"github.com/tanxiaolong/apidoc/locale"
	"github.com/tanxiaolong/apidoc/output/static"
	"github.com/tanxiaolong/apidoc/types"
	"github.com/tanxiaolong/apidoc/vars"

	"github.com/issue9/utils"
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	Dir     string        `yaml:"dir"`              // 文档的保存目录
	Groups  []string      `yaml:"groups,omitempty"` // 仅输出这些组，为空表示输出所有
	Elapsed time.Duration `yaml:"-"`                // 编译用时

	dataDir string // json 数据保存的目录
}

// Sanitize 对 Options 作一些初始化操作。
func (o *Options) Sanitize() *types.OptionsError {
	if len(o.Dir) == 0 {
		return &types.OptionsError{Field: "dir", Message: locale.Sprintf(locale.ErrRequired)}
	}

	return nil
}

// 指定的组是否需要输出
func (o *Options) groupIsEnable(group string) bool {
	if len(o.Groups) == 0 {
		return true
	}

	for _, g := range o.Groups {
		if g == group {
			return true
		}
	}

	return false
}

// Render 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(docs *types.Doc, o *Options) error {
	// 文档目录下的文件名可能改变，先清除目录下的所有文件。
	if err := os.RemoveAll(o.Dir); err != nil {
		return err
	}

	if !utils.FileExists(o.Dir) {
		if err := os.MkdirAll(o.Dir, os.ModePerm); err != nil {
			return err
		}
	}

	o.dataDir = filepath.Join(o.Dir, vars.JSONDataDirName)
	if !utils.FileExists(o.dataDir) {
		if err := os.MkdirAll(o.dataDir, os.ModePerm); err != nil {
			return err
		}
	}

	if err := static.Output(o.Dir); err != nil {
		return err
	}

	return render(docs, o)
}
