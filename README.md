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