// Author: peng.lee577 (peng.lee577@gmail.com)
// Date: 2023/9/13

package env

import (
	"fmt"
	"path/filepath"
)

var myEnv AppEnv

func init() {
	myEnv = New(Option{
		RootDir: Rootdir,
	})
}

func ExampleRootDir() {
	dir := myEnv.RootDir()
	fmt.Println("rootDir=", dir)
	// Output:
	// rootDir= ./testdata/
}

func ExampleConfDir() {
	dir := myEnv.ConfDir()
	fmt.Println("confRootDir=", dir)

	dbUserConfPath := filepath.Join(myEnv.ConfDir(), "db", "db_user.toml")
	fmt.Println("dbConfPath=", dbUserConfPath)

	// Output:
	// confRootDir= testdata/conf
	// dbConfPath= testdata/conf/db/db_user.toml
}

func ExampleDataDir() {
	dir := myEnv.DataDir()
	fmt.Println("dataRootDir=", dir)

	// Output:
	// dataRootDir= testdata/data
}

func ExampleLogDir() {
	dir := myEnv.LogDir()
	fmt.Println("logRootDir=", dir)

	dbLogPath := filepath.Join(dir, "db", "mysql.log")
	fmt.Println("dbLogPath=", dbLogPath)

	// Output:
	// logRootDir= testdata/log
	// dbLogPath= testdata/log/db/mysql.log
}

func ExampleNew() {
	// 比如在应用初始化的阶段，可以通过如下方式来修改默认的环境信息
	opt := Option{
		AppName: "demo",
		IDC:     "jx",
		RunMode: RunModeDebug,
		RootDir: Rootdir,

		// 下面这几项可选
		DataDir: Rootdir + "my_data",
		LogDir:  Rootdir + "my_log",
		ConfDir: Rootdir + "my_conf",
	}

	Default = New(opt)
}

func ExampleRunMode() {
	switch RunMode() {
	case RunModeDebug:
		fmt.Println("debug")

	case RunModeTest:
		fmt.Println("test")

	case RunModeRelease:
		fmt.Println("release")
	}
}

func ExampleIDC() {
	if IDC() == "jx" {
		fmt.Println("jx")
	}
}
