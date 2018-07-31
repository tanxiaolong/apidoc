# 如何为 apidoc 贡献代码

apidoc 是一个基于 [MIT](https://opensource.org/licenses/MIT) 的开源软件。
欢迎大家共同参与开发。**若需要新功能，请先开 issue 讨论。**



## 本地化

本地化的内容在 `locale` 包中，欢迎大家对即有的内容作出翻译修改，
但暂时不支持添加新语言的支持。



## 添加新语言

`input/lang.go` 文件中有所有语言模型的定义，若需要添加对新语言的支持，
在该文件中有详细的文档说明如何定义语言模弄。提交后请更新
[#11](https://github.com/tanxiaolong/apidoc/issues/11)



## 模板

默认模板在 `output/static` 目录下，目前要求支持的浏览器为：
Edge、Chrome、Firefox 和 Safari 的最后一个稳定版本。



## 文档

要求支持的浏览器为：Edge、Chrome、Firefox 和 Safari 的最后一个稳定版本。
