# TypePress

TypePress 是一个 Blog 站群系统, 以 [go][0] 语言实现. 初衷是学习和实践 [go][0] 语言下的 WEB 开发. 实践不同开发方法对后续维护的影响. 当然作为一个 Blog 系统, 实用性和易用性是必须要考虑的.

## 敏感资料加密

Blog 系统是需要用户注册才能使用. TypePress 对敏感的基本资料, 比如登录名和密码进行了加密(MD5)存储. 并且这种加密是在浏览器中进行的. 也就是说正常情况下, 系统无法获取用户登录名和密码的原始值. 连系统都不知道, 自然无法泄密. 非正常情况, 比如用 email 找回密码, 才有可能让系统知道用户的真正 email. 当然如果用户愿意公开这些信息, 用户是有这个选择的.

作为开源软件, TypePress 无法控制使用者破坏这种保护措施.

TypePress 提醒最终用户, 使用 TypePress 且不遵守这种保护措施的站点, 属于不良设计, 怀疑有非善意目的.

可以在浏览器上监视到网站是否遵守这种保护. 

## 框架

框架是客观存在的, TypePress 更关心通过组合独立 package 来完成任务, 而不是提供或者使用一个大而全的框架. 当然实现这个设想是比较困难的, 某些地方很难区分是否够"独立"并解耦, 是否已经算是框架了.
作为尝试, 设计初期 TypePress 不知道会遇到什么情况, 这种想法彻底失败也有可能.

## 开发过程

整个开发过程在 [Go-Blog-In-Action][6].

## 使用

### 获取源码

以我使用机器环境变量

>GOPATH=F:\go

>GOROOT=E:\Go

为例. 使用

```
go get github.com/achun/typepress
```

会得到这样的提示

```
package github.com/achun/typepress
        imports github.com/achun/typepress
        imports github.com/achun/typepress: no Go source files in F:\go\src\github.com\achun\typepress
```

出现 `no Go source files` 是正常的, 因为 TypePress 目录中增加了 `src` 子目录. 其实只是下载了源码, 相关依赖 package 并未得到自动安装.

### 配置开发环境变量

如果使用 `Sublime Text` + `GoSublime`, 可以通过菜单

```
Preferences -> Package Settings -> GoSublime -> Settings - User
```

设置

```
"env": { "GOPATH": "$GOPATH;$GS_GOPATH" }
```

如果使用 LiteIDE, 可以通过编辑环境变量的方法给相应 `.env` 配置添加绝对路径

```
GOPATH=F:\go;F:\go\src\github.com\achun\typepress
```

也可以通过菜单

```
查看->管理GOPATH
```

进行添加.

建议不要把 TypePress 路径添加到系统GOPATH中.

### 直接修改 TypePress 源码进行使用

这是最简单的一种使用方法, 当然如果使用这种方法, 复制一份 TypePress 的拷贝进行修改是个好方法.

### 通过 import 方式使用

上面已经把 TypePress 路径已经加入到 `GOPATH` (开发环境下)中, 所以在您的项目中

```go
import "global"
```

这种用法完全没有问题, 您也可以用修改环境变量的方法, 把您的项目路径 `$YourPackAgePath` 加入到开发环境的 `GOPATH` 中

```
"env": { "GOPATH": "$GOPATH;$YourPackAgePath;$GS_GOPATH" }
```
或者
```
GOPATH=F:\go;$YourPackAgePath;F:\go\src\github.com\achun\typepress
```

### 需要手工 go get 的 package

例如, 经过上面的设置后, 在 Sublime Text 中把 TypePress 目录加入 FOLDERS 并 `go build` `main.go`.
有可能会遇到
```
cannot find package "github.com/achun/db" in any of: ...
```
这样的错误, 原因就是前面所述 TypePress 所依赖的 package 没有被自动安装造成的.

您需要手工 go get 那些缺失的 package. 下面的代码方便您使用, 也许还有其他.
```
go get -u github.com/achun/template
go get -u github.com/achun/log
go get -u github.com/achun/db
go get -u github.com/achun/go-toml
go get -u github.com/gorilla/context
go get -u github.com/gorilla/mux
go get -u github.com/braintree/manners
```

## License

TypePress 采用
MIT License: http://achun.mit-license.org

TypePress 只使用采用下列 License 的 Repository.

* [MIT][1]
* [BSD-2-Clause][2] 
* [BSD-3-Clause][3] 
* [Apache v2 License][4]
* [Public Domain Unlicense][5]

[0]: https://golang.org
[1]: http://choosealicense.com/licenses/mit/
[2]: http://choosealicense.com/licenses/bsd/
[3]: http://choosealicense.com/licenses/bsd-3-clause/
[4]: http://choosealicense.com/licenses/apache/
[5]: http://choosealicense.com/licenses/public-domain/
[6]: https://github.com/achun/Go-Blog-In-Action