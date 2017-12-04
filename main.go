// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/output"
	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"

	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	tag, err := locale.Init()
	if err != nil {
		warn.Println(err)
		info.Println("无法获取系统语言，使用默认的本化语言：", vars.DefaultLocale)
		tag, err = language.Parse(vars.DefaultLocale)
		if err != nil {
			erro.Println(err)
			return
		}
	}
	locale.SetLocale(tag)

	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	l := flag.Bool("l", false, locale.Sprintf(locale.FlagLUsage))
	g := flag.Bool("g", false, locale.Sprintf(locale.FlagGUsage))
	pprofType := flag.String("pprof", "", locale.Sprintf(locale.FlagPprofUsage))
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		printVersion()
		return
	case *l:
		locale.Printf(locale.FlagSupportedLangs, input.Langs())
		return
	case *g:
		path, err := getConfigFile()
		if err != nil {
			erro.Println(err)
			return
		}
		if err = genConfigFile(path); err != nil {
			erro.Println(err)
			return
		}
		info.Println(locale.Sprintf(locale.FlagConfigWritedSuccess, path))
		return
	}

	// 指定了 pprof 参数
	if len(*pprofType) > 0 {
		profile := filepath.Join("./", vars.Profile)
		f, err := os.Create(profile)
		if err != nil { // 不能创建文件，则忽略 pprof 相关操作
			warn.Println(err)
			goto RUN
		}
		defer func() {
			if err = f.Close(); err != nil {
				erro.Println(err)
				return
			}
			info.Println(locale.Sprintf(locale.FlagPprofWritedSuccess, profile))
		}()

		switch strings.ToLower(*pprofType) {
		case "mem":
			defer func() {
				if err = pprof.Lookup("heap").WriteTo(f, 1); err != nil {
					warn.Println(err)
				}
			}()
		case "cpu":
			if err := pprof.StartCPUProfile(f); err != nil {
				warn.Println(err)
			}
			defer pprof.StopCPUProfile()
		default:
			erro.Println(locale.Sprintf(locale.FlagInvalidPprrof))
			return
		}
	}

RUN:
	run()
}

// 真正的程序入口，main 主要是作参数的处理。
func run() {
	start := time.Now()

	path, err := getConfigFile()
	if err != nil {
		erro.Println(err)
		return
	}

	cfg, err := loadConfig(path)
	if err != nil {
		erro.Println(err)
		return
	}

	// 分析文档内容
	docs := types.NewDoc()
	wg := &sync.WaitGroup{}
	for _, opt := range cfg.Inputs {
		wg.Add(1)
		go func(o *input.Options) {
			if err := input.Parse(docs, o); err != nil {
				erro.Println(err)
			}
			wg.Done()
		}(opt)
	}
	wg.Wait()

	if len(docs.Title) == 0 {
		docs.Title = vars.DefaultTitle
	}

	// 输出内容
	cfg.Output.Elapsed = time.Now().Sub(start)
	if err := output.Render(docs, cfg.Output); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.Complete, cfg.Output.Dir, time.Now().Sub(start)))
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	locale.Printf(locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

// 获取配置文件路径。目前只支持从工作路径获取。
func getConfigFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, vars.ConfigFilename), nil
}

// 生成一个默认的配置文件，并写入到 path 中。
func genConfigFile(path string) error {
	dir := filepath.Dir(path)
	lang, err := input.DetectDirLang(dir)
	if err != nil { // 不中断，仅作提示用。
		warn.Println(err)
	}

	cfg := &config{
		Version: vars.Version(),
		Inputs: []*input.Options{
			&input.Options{
				Dir:       dir,
				Recursive: true,
				Lang:      lang,
			},
		},
		Output: &output.Options{
			Dir: filepath.Join(dir, "doc"),
		},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Write(data)
	return err
}

func printVersion() {
	locale.Printf(locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version())
	locale.Printf(locale.FlagVersionCommitHash, vars.CommitHash())
}
