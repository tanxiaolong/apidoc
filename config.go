// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"strconv"

	"github.com/issue9/version"
	yaml "gopkg.in/yaml.v2"

	"github.com/tanxiaolong/apidoc/input"
	"github.com/tanxiaolong/apidoc/locale"
	"github.com/tanxiaolong/apidoc/output"
	"github.com/tanxiaolong/apidoc/types"
)

// 项目的配置内容，分别引用到了 input.Options 和 output.Options。
//
// 所有可能改变输出的表现形式的，应该添加此配置中；
// 而如果只是改变输出内容的，应该直接以标签的形式出现在代码中，
// 比如文档的版本号、标题等，都是直接使用 @apidoc 来指定的。
type config struct {
	Version string           `yaml:"version"` // 产生此配置文件的程序版本号
	Inputs  []*input.Options `yaml:"inputs"`  // 输入的配置项，可以指定多个项目
	Output  *output.Options  `yaml:"output"`
}

// 加载 path 所指的文件内容到 *config 实例。
func loadConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err = cfg.sanitize(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *config) sanitize() error {
	if !version.SemVerValid(cfg.Version) {
		return &types.OptionsError{Field: "version", Message: locale.Sprintf(locale.ErrInvalidFormat)}
	}

	if len(cfg.Inputs) == 0 {
		return &types.OptionsError{Field: "inputs", Message: locale.Sprintf(locale.ErrRequired)}
	}

	if cfg.Output == nil {
		return &types.OptionsError{Field: "output", Message: locale.Sprintf(locale.ErrRequired)}
	}

	for i, opt := range cfg.Inputs {
		if err := opt.Sanitize(); err != nil {
			index := strconv.Itoa(i)
			err.Field = "inputs[" + index + "]." + err.Field
			return err
		}

		// 指定语法检测输出通道
		opt.ErrorLog = erro
		opt.WarnLog = warn
	}

	if err := cfg.Output.Sanitize(); err != nil {
		err.Field = "outputs." + err.Field
		return err
	}

	return nil
}
