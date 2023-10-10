// Author: peng.lee577 (peng.lee577@gmail.com)
// Date: 2020/7/13

package env

import (
	"os"
	"strings"
)

// AttributeReader 读取属性的接口定义
// 为了方便序列化、以及和系统环境变量打通，以及最终的使用，属性的 key 和 value 都定义为 string 类型
// 而不是 any 类型
//
// 此处定义的属性和 其他的如 IDC、APPName 等值是独立的，功能是互不影响的
type AttributeReader interface {
	// Attribute 读取一个属性，若是不存在，会返回 nil
	// 若是从存储中没有查询到，同时 key 是 string 类型，会尝试从环境变量中查询：
	// 如 查询 key = "abc",存储中没有，会尝试从环境变量 "GDP_ENV_abc" 中读取
	Attribute(key string) string
}

// AttributeWriter 修改属性的接口定义
type AttributeWriter interface {
	// SetAttribute 需要支持运行时动态更新
	// 若 val == ""，则删除相应的 key
	SetAttribute(key string, val string)
}

// 特殊的属性读取方法
var specAttrsHooks = map[string]func() string{
	"app": attrAppNameFromOsEnv,
}

// 通过环境变量获取应用名称
// Matrix 环境：
// CONTAINER_ID = 1000004.bdapp-gdp-website-tucheng
// 最终返回 bdapp-gdp-website-tucheng
//
// 若是在 EKS 上（无有效 CONTAINER_ID 环境变量）：
// 将返回 $APPSPACE_NAMESPACE+"-"+ $APPSPACE_IDC_NAME
func attrAppNameFromOsEnv() string {
	// Matrix 环境：
	// CONTAINER_ID = 1000004.bdapp-gdp-website-tucheng
	id := os.Getenv("CONTAINER_ID")
	start := strings.Index(id, ".")
	if start > 0 {
		// 返回 bdapp-gdp-website-tucheng
		return id[start+1:]
	}
	// EKS
	eksName := os.Getenv("APPSPACE_NAMESPACE")
	eksIDC := os.Getenv("APPSPACE_IDC_NAME")
	if len(eksIDC) > 0 && len(eksName) > 0 {
		return eksName + "-" + eksIDC
	}
	return ""
}
