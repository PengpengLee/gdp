// Author: peng.lee577 (peng.lee577@gmail.com)
// Date: 2023/7/13

package env

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var Rootdir = "./testdata/"

func TestAutoDetectAppRootDir(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "case 1",
			want: wd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AutoDetectAppRootDir()
			if got != tt.want {
				t.Errorf("AutoDetectAppRootDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findDirMatch(t *testing.T) {
	type args struct {
		baseDir   string
		fileNames []string
	}
	tests := []struct {
		name    string
		args    args
		wantDir string
		wantErr bool
	}{
		{
			name: "dir1_a.txt",
			args: args{
				baseDir:   Rootdir + "dir1",
				fileNames: []string{"a.text"},
			},
			wantDir: Rootdir + "dir1",
			wantErr: false,
		},
		{
			name: "miss",
			args: args{
				baseDir:   Rootdir + "dir1",
				fileNames: []string{"a.not_exists"},
			},
			wantDir: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDir, err := findDirMatch(tt.args.baseDir, tt.args.fileNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("findDirMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDir != tt.wantDir {
				t.Errorf("findDirMatch() gotDir = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func TestNewAppEnv(t *testing.T) {
	wd, _ := os.Getwd()
	root := wd

	tests := []struct {
		name      string
		argValue  Option
		wantValue Option
	}{
		{
			name:     "case 1",
			argValue: Option{},
			wantValue: Option{
				RootDir: root,
				DataDir: filepath.Join(root, "data"),
				LogDir:  filepath.Join(root, "log"),
				ConfDir: filepath.Join(root, "conf"),
				IDC:     "test",
				AppName: "unknown",
				RunMode: "release",
			},
		},
		{
			name: "case 2",
			argValue: Option{
				RootDir: Rootdir,
				ConfDir: Rootdir + "conf1",
				LogDir:  Rootdir + "log1",
				DataDir: Rootdir + "data1",
				IDC:     "jx",
				AppName: "abc",
				RunMode: "debug",
			},
			wantValue: Option{
				RootDir: Rootdir,
				ConfDir: Rootdir + "conf1",
				LogDir:  Rootdir + "log1",
				DataDir: Rootdir + "data1",
				IDC:     "jx",
				AppName: "abc",
				RunMode: "debug",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.argValue
			env := New(val)

			val.LogDir = "不会修改环境信息"

			got := env.Options()
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("got=%s,\n want=%s", got, tt.wantValue)
			}
		})
	}
}

func Test_setValue(t *testing.T) {
	var a string
	setValue(&a, "ok", "set a")
	if a != "ok" {
		t.Errorf("a=%q, want=%q", a, "ok")
	}
	func() {
		defer func() {
			if re := recover(); re != nil {
				t.Error("expect no panic")
			}
		}()
		setValue(&a, "panic", "set a")
	}()
}

func TestValue_String(t *testing.T) {
	val := Option{
		RootDir: "dir_a",
		DataDir: "dir_b",
		LogDir:  "dir_c",
		ConfDir: "dir_d",
	}
	got := val.String()
	var v1 *Option
	if err := json.Unmarshal([]byte(got), &v1); err != nil {
		t.Error("json.Unmarshal with error:", err)
	}

	if !reflect.DeepEqual(v1, &val) {
		t.Errorf("Unmarshal got=%v, want=%v", v1, val)
	}
}

func TestOption_Merge(t *testing.T) {
	tests := []struct {
		name string
		opt  Option
		args Option
		want Option
	}{
		{
			name: "case 1",
			opt: Option{
				AppName: "a",
				IDC:     "b",
				RunMode: "c",
				RootDir: "d",
				DataDir: "e",
				LogDir:  "f",
				ConfDir: "g",
			},
			args: Option{},
			want: Option{
				AppName: "a",
				IDC:     "b",
				RunMode: "c",
				RootDir: "d",
				DataDir: "e",
				LogDir:  "f",
				ConfDir: "g",
			},
		},
		{
			name: "case 2",
			opt: Option{
				AppName: "a",
				IDC:     "b",
				RunMode: "c",
				RootDir: "d",
				DataDir: "e",
				LogDir:  "f",
				ConfDir: "g",
			},
			args: Option{
				AppName: "1",
				IDC:     "2",
				RunMode: "3",
				RootDir: "4",
				DataDir: "5",
				LogDir:  "6",
				ConfDir: "7",
			},
			want: Option{
				AppName: "1",
				IDC:     "2",
				RunMode: "3",
				RootDir: "4",
				DataDir: "5",
				LogDir:  "6",
				ConfDir: "7",
			},
		},
		{
			name: "case 3",
			opt: Option{
				AppName: "a",
				IDC:     "b",
				RunMode: "c",
				RootDir: "d",
				DataDir: "e",
				LogDir:  "f",
				ConfDir: "g",
			},
			args: Option{
				AppName: "1",
				IDC:     "2",
				RunMode: "3",
			},
			want: Option{
				AppName: "1",
				IDC:     "2",
				RunMode: "3",
				RootDir: "d",
				DataDir: "e",
				LogDir:  "f",
				ConfDir: "g",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.opt
			if got := opt.Merge(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
