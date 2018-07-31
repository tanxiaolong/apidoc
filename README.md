apidoc [![Build Status](https://travis-ci.org/tanxiaolong/apidoc.svg?branch=master)](https://travis-ci.org/tanxiaolong/apidoc)
[![Go version](https://img.shields.io/badge/Go-1.8-brightgreen.svg?style=flat)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/tanxiaolong/apidoc)](https://goreportcard.com/report/github.com/tanxiaolong/apidoc)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/tanxiaolong/apidoc/branch/master/graph/badge.svg)](https://codecov.io/gh/tanxiaolong/apidoc)
======

apidoc 是一个简单的 RESTful API 文档生成工具，它从代码注释中提取特定格式的内容，生成文档。
目前支持支持以下语言：C#、C/C++、D、Erlang、Go、Groovy、Java、JavaScript、Pascal/Delphi、
Perl、PHP、Python、Ruby、Rust、Scala 和 Swift。

具体文档可参考：http://apidoc.tools

```go
/**
 * @api get /users 获取所有的用户信息
 * @apiGroup users
 * @apiQuery page int 显示第几页的内容
 * @apiQuery size int 每页显示的数量
 *
 * @apiSuccess 200 ok
 * @apiParam count int 符合条件的所有用户数量
 * @apiParam users array 用户列表。
 * @apiExample json
 * {
 *     "count": 500,
 *     "users": [
 *         {"id":1, "username": "admin1", "name": "管理员2"},
 *         {"id":2, "username": "admin2", "name": "管理员2"}
 *     ],
 * }
 * @apiExample xml
 * <users count="500">
 *     <user id="1" username="admin1" name="管理员1" />
 *     <user id="2" username="admin2" name="管理员2" />
 * </users>
 */
func login(w http.ResponseWriter, r *http.Request) {
    // TODO
}
```



### 安装

```shell
go get github.com/tanxiaolong/apidoc
```

支持多种本地化语言，默认情况下会根据当前系统所使用的语言进行调整。若需要手动指定，
windows 可以设置一个 `LANG` 环境变量指定，*nix 系统可以使用以下命令：
```shell
LANG=lang apidoc
```
将其中的 lang 设置为你需要的语言。



### 集成

若需要将 apidoc 当作包集成到其它 Go 程序中，可分别引用 `input` 和 `output` 的相关函数：

```go
// 初始本地化内容
locale.Init()

// 分析文档内容
inputOptions := &input.Options{
    ErrorLog: log.New(...),
}
docs, elapsed := input.Parse(inputOptions)

// 输出内容
outputOptions := &output.Options{...}
outputOptions.Elapsed = elapsed
if err = output.Render(docs, outputOptions); err != nil {
    panic(err)
}
```



### 参与开发

请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 文件的相关内容。

### 配置文件示例
```go
version: 4.0.0+20171125
inputs:
    - lang: go
      dir: /home/tanxiaolong/code/
      recursive: true
      title: 用户中心-互动接口文档
output:
    dir: /www/vcs-doc
    type: html
```

### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
