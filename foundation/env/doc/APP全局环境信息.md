# APP全局环境信息



# 1.概述

用于管理全局的应用信息。

自动推断出项目根目录、配置目录、数据目录等。



# 2.接口设计

## 2.1 基础定义和实现

```go
// AppEnv 应用环境信息完整的接口定义
type AppEnv interface {
   // 应用名称
   AppNameEnv

   // idc信息
   IDCEnv

   // 应用根目录
   RootDirEnv

   // 应用配置文件根目录
   ConfDirEnv

   // 应用数据文件根目录
   DataDirEnv

   // 应用日志文件更目录
   LogDirEnv

   // 应用运行情况
   RunModeEnv

   // 获取当前环境的选项详情
   Options() Option

   // 复制一个新的env对象，并将传入的Option merge进去
   CloneWithOption(opt Option) AppEnv
   
   AttributeReader
   AttributeWriter
}

// RootDirEnv 应用根目录环境信息
type RootDirEnv interface {
   RootDir() string
}

// ConfDirEnv 配置环境信息
type ConfDirEnv interface {
   ConfDir() string
}

// DataDirEnv 数据目录环境信息
type DataDirEnv interface {
   DataDir() string
}

// LogDirEnv 日志目录环境信息
type LogDirEnv interface {
   LogDir() string
}

// IDCEnv idc的信息读写接口
type IDCEnv interface {
   IDC() string
}

// AppNameEnv 应用名称
type AppNameEnv interface {
   AppName() string
}

// RunModeEnv 运行模式/等级
type RunModeEnv interface {
   RunMode() string
}
```



## 2.2 全局方法

```go
// Default (全局)默认的环境信息
//
// 全局的 RootDir() 、DataDir() 等方法均使用该环境信息
var Default = New(Option{})

// RootDir (全局)获取应用根目录
func RootDir() string {
   return Default.RootDir()
}

// DataDir (全局)设置应用数据根目录
func DataDir() string {
   return Default.DataDir()
}

// LogDir (全局)获取应用日志根目录
func LogDir() string {
   return Default.LogDir()
}

// ConfDir (全局)获取应用配置根目录
func ConfDir() string {
   return Default.ConfDir()
}

// IDC (全局)获取应用的idc
func IDC() string {
   return Default.IDC()
}

// AppName (全局)应用的名称
func AppName() string {
   return Default.AppName()
}

// RunMode (全局) 程序运行等级
// 默认是 release(线上发布)，还可选 RunModeDebug、RunModeTest
func RunMode() string {
   return Default.RunMode()
}

// Options 获取当前环境的选项详情
func Options() Option {
   return Default.Options()
}

// CloneWithOption 复制一个新的env对象，并将传入的Option merge进去
func CloneWithOption(opt Option) AppEnv {
   return Default.CloneWithOption(opt)
}
```



## 2.3 扩展-Attributes 功能

**定义说明：**

   Attributes 是扩展属性(Extension Attributes)，是用来给应用添加自定义 Label/属性的，这样跨模块的属性传递会更加方便。

   如在业务模块中定义特定的 Label，框架的服务发现模块可以读取并使用这些 Label，以此做出相应决策。如在[匹配条件-sourceLabels](#1. sourceLabels (持续完善))中有说明相关的使用。

```go
// AttributeReader 读取属性的接口定义
// 为了方便序列化、以及和系统环境变量打通，以及最终的使用，属性的 key 和 value 都定义为 string 类型
// 而不是 any 类型
//
// 此处定义的属性和 其他的如 IDC、APPName 等值是独立的，功能是互不影响的
type AttributeReader interface {
	// Attribute 读取一个属性，若是不存在，会返回 nil
	Attribute(key string) string
}

// AttributeWriter 修改属性的接口定义
type AttributeWriter interface {
	// SetAttribute 需要支持运行时动态更新
	// 若 val == ""，则删除相应的 key
	SetAttribute(key string, val string)
}
```

原来的 `AppEnv` 接口新增 `AttributeReader` 和 `AttributeWriter` 约束。同时新增全局方法 `Attribute(key string) string` 和 `SetAttribute(key string, val string)` 以方便使用。

属性读取的时候，若读取不到，会尝试从环境变量中读取（能简化使用，利用应用已有的环境变量属性，方便程序在运行过程中动态更新策略），具体策略定义如下：

```go
// Attribute 读取一个属性，若是不存在，会返回 nil
// Default.Attribute 的别名
// 若是从存储中没有查询到，同时 key 是 string 类型，会尝试从环境变量中查询：
// 如 查询 key = "abc",存储中没有，会尝试从环境变量 "GDP_ENV_abc" 中读取
//
// 注意：IDC 等其他环境信息应使用该 pkg 的方法读取，如 IDC() 方法
// 若 SetAttribute("IDC","jx") ,并不会影响 IDC() 方法的返回值
//
// 若读取 key="app",当存储中没有，同时环境变量 GDP_ENV_app 也不存在的时候
// 会先尝试从环境变量 CONTAINER_ID 中解析：
//
//	该环境变量要求格式为：CONTAINER_ID = 1000004.bdapp-gdp-website-tucheng
//	目前厂内的各种基于 Matrix 的平台都有该环境变量
//
// 若是在 EKS 上（无有效 CONTAINER_ID 环境变量）:
//
//	将返回 $APPSPACE_NAMESPACE+"-"+ $APPSPACE_IDC_NAME
//
// 若 key 是以 "OE_" 为前缀，则最终尝试从系统环境变量中读取
// 如 key="OE_PIDC",则从系统环境变量 PIDC 读取
func Attribute(key string) string {
	return Default.Attribute(key)
}

// SetAttribute 支持运行时动态更新
// Default.SetAttribute 的别名
// 若 val == ""，则删除相应的 key
func SetAttribute(key string, val string) {
	Default.SetAttribute(key, val)
}
```

# 3.部分细节

## 3.1 应用根目录自动推断

**需求背景：**

应用程序打印日志，读取配置等都涉及到程序根目录的识别，另外在运行单元测试的时候，由于可能是直接在源码所在目录下运行的。

在v1版本里，运行 `go test` 的时候，就有可能不能正确的找到程序的根目录，导致配置不能正确的读取到，日志打印的目录是错误的，这也就导致了写单测的成本加大。

新版本的自动推断应用根目录不正确的问题需要去解决这个case。

**实现方案：**

由于是要解决在运行 `go test` 过程的时候根目录的识别。所以可以利用代码的目录结构来进行识别。

```cmd
├── gconn
│   ├── gconn.go
│   └── tracer.go
├── go.mod
├── go.sum
├── http
│   ├── request_response.go
│   ├── request_response_test.go
│   ├── tracer.go
│   └── tracer_test.go
├── logit
│   ├── callstack.go
│   ├── codegen.py
│   ├── context.go
│   ├── duration.go
│   ├── encoder.go
│   ├── field.go
│   ├── field_test.go
│   ├── gdpbridge.go
│   ├── gdpbridge_test.go
│   ├── level.go
│   ├── logger.go
│   ├── noplogger.go
│   ├── noplogger_test.go
│   ├── simplelogger.go
│   └── simplelogger_test.go
├── main.go
```

如上是一个项目代码结构，可以看到 go.mod 文件总是在项目根目录下的，所以，可以利用向上查找go.mod 文件来定位程序根目录。

##  3.2 系统环境变量

注：此部分为 2022 年 11 月 28 日新增

此模块从系统环境变量读取的 key 的前缀为 `GDP_ENV_`

如下默认值会优先使用系统环境变量的值：

```go
var (
	// DefaultIDC 默认 idc 的值
	// idc 推荐可选值  如 jx,gz等
	// 默认值会优先使用环境变量 'GDP_ENV_IDC' 的值，若没有，则使用 test
	DefaultIDC = osEnvDefault("IDC", "test")

	// DefaultAppName 默认的 app 名称
	// 默认值会优先使用环境变量 'GDP_ENV_AppName' 的值，若没有，则使用 unknown
	DefaultAppName = osEnvDefault("AppName", "unknown")

	// DefaultRunMode 测试默认运行等级
	// 默认值会优先使用环境变量 'GDP_ENV_RunMode' 的值，若没有，则使用 RunModeRelease
	DefaultRunMode = readOsEnvRunMode()
)
```

另外 上述 Attributes 当 key 没有读取到值的时候，也会重试从系统变量中读取。



# 附录：

## 1. sourceLabels (持续完善)

sourceLabels，用于限定客户端的条件，可以有 0 或者多个。社区最常用的 label 是 "app" 和 “version” 这两个字段。

本地客户端也会有对应的值。

**sourceLabels 使用示例：**

| 匹配条件                     | 客户端的配置             | 是否匹配 |
| ---------------------------- | ------------------------ | -------- |
| app="demo" <br>version="v1"  | app="demo"               | 是       |
| app="demo" <br/>version="v1" | app="demo" <br/>app="v2" | 否       |