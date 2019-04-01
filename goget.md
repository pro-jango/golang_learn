### go get命令——一键获取代码、编译并安装
#### 命令说明
go get 可以借助代码管理工具通过远程拉取或更新代码包及其依赖包，并自动完成编译和安装。整个过程就像安装一个 App 一样简单。

使用 go get 前，需要安装与远程包匹配的代码管理工具，如 Git、SVN、HG 等，参数中需要提供一个包名。 

```
$ go get github.com/xxx/xxxx

```


#### go get使用时的附加参数
附加参数|备注
---|---
-v|显示操作流程的日志及信息，方便检查错误
-u|下载丢失的包，但不会更新已经存在的包
-d|只下载，不安装
-insecure|允许使用不安全的 HTTP 方式进行下载操作

```
go get -v github.com/xxx/xxx
go get -d github.com/xxx/xxx
go get -u github.com/xxx/xxx
```