package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error:%v", err.Error())
	}

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
			if got := RootDir(); got != tt.want {
				t.Errorf("Rootdir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testResetDefault() {
	Default = New(Option{})
}

func TestLogRootDir(t *testing.T) {
	defer testResetDefault()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error:%v", err.Error())
	}
	tests := []struct {
		name string
		call func()
		want string
	}{
		{
			name: "case 1",
			want: filepath.Join(wd, "log"),
		},
		{
			name: "case 2",
			call: func() {
				Default = New(Option{LogDir: "testdata"})
			},
			want: "testdata",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.call != nil {
				tt.call()
			}
			if got := LogDir(); got != tt.want {
				t.Errorf("LogDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfRootPath(t *testing.T) {
	defer testResetDefault()

	tests := []struct {
		name string
		call func()
		want string
	}{
		{
			name: "case 1",
			call: func() {
				Default = New(Option{RootDir: "internal"})
			},
			want: "internal/conf",
		},
		{
			name: "case 2",
			call: func() {
				Default = New(Option{
					ConfDir: "xyz/abc",
				})
			},
			want: "xyz/abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.call()
			got := ConfDir()
			if got != tt.want {
				t.Errorf("got=%q, want=%q", got, tt.want)
			}
		})
	}

	// 恢复初始化的环境
	Default = New(Option{})
}

func TestDataRootPath(t *testing.T) {
	defer testResetDefault()
	tests := []struct {
		name string
		call func()
		want string
	}{
		{
			name: "case 1",
			call: func() {
				Default = New(Option{RootDir: "./internal"})
			},
			want: "internal/data",
		},
		{
			name: "case 2",
			call: func() {
				Default = New(Option{
					DataDir: "xyz/abc",
				})
			},
			want: "xyz/abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.call()
			got := DataDir()
			if got != tt.want {
				t.Errorf("got=%q, want=%q", got, tt.want)
			}
		})
	}

	// 恢复初始化的环境
	Default = New(Option{})
}

func TestIDC(t *testing.T) {
	defer testResetDefault()

	if got := IDC(); got != "test" {
		t.Errorf("IDC()=%q want=%q", got, "test")
	}

	Default = New(Option{IDC: "jx"})

	if got := IDC(); got != "jx" {
		t.Errorf("IDC()=%q want=%q", got, "jx")
	}

	opts := Options()
	if opts.IDC != "jx" {
		t.Errorf("opts.IDC=%q want=%q", opts.IDC, "jx")
	}
}

func TestAppName(t *testing.T) {
	defer testResetDefault()

	if got := AppName(); got != "unknown" {
		t.Errorf("AppName()=%q want=%q", got, "unknown")
	}

	Default = New(Option{AppName: "demo"})

	if got := AppName(); got != "demo" {
		t.Errorf("AppName()=%q want=%q", got, "demo")
	}

	b := Default.CloneWithOption(Option{})
	if got := b.AppName(); got != "demo" {
		t.Errorf("b.AppName()=%q want=%q", got, "demo")
	}

	c := CloneWithOption(Option{AppName: "test"})
	if got := c.AppName(); got != "test" {
		t.Errorf("c.AppName()=%q want=%q", got, "test")
	}
}

func TestRunMode(t *testing.T) {
	defer func() {
		Default = New(Option{})
	}()
	if got := RunMode(); got != "release" {
		t.Errorf("RunMode()=%q want=%q", got, "release")
	}

	Default = New(Option{RunMode: "test"})

	if got := RunMode(); got != "test" {
		t.Errorf("RunMode()=%q want=%q", got, "test")
	}
}
