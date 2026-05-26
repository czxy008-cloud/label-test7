//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/magefile/mage/sh"
)

const (
	appName     = "fresh-groupbuy"
	mainPath  = "./cmd/main.go"
	outputDir = "./bin"
)

var (
	goexe = sh.RunCmd("go")

func init() {
	os.Setenv("CGO_ENABLED", "0")
}

func Build() error {
	fmt.Println("开始构建...")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	version := getVersion()
	commit := getCommit()
	buildTime := time.Now().Format("2006-01-02 15:04:05")

	ldflags := fmt.Sprintf(
		"-s -w -X main.Version=%s -X main.Commit=%s -X main.BuildTime=%q",
		version, commit, buildTime,
	)

	output := filepath.Join(outputDir, appName)
	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	fmt.Printf("构建目标: %s", output)

	return goexe("build",
		"-ldflags", ldflags,
		"-o", output,
		mainPath,
	)
}

func Run() error {
	if err := Build(); err != nil {
		return err
	}
	output := filepath.Join(outputDir, appName)
	if runtime.GOOS == "windows" {
		output += ".exe"
	}
	return sh.RunV(output)
}

func Test() error {
	fmt.Println("运行测试...")
	return goexe("test", "./...", "-v", "-cover")
}

func Lint() error {
	fmt.Println("运行代码检查...")
	return sh.RunV("golangci-lint", "run", "./...")
}

func Clean() error {
	fmt.Println("清理构建产物...")
	if _, err := os.Stat(outputDir); err == nil {
		return os.RemoveAll(outputDir)
	}
	return nil
}

func InitDB() error {
	fmt.Println("初始化数据库...")
	return sh.RunV("psql", "-U", "postgres", "-f", "./scripts/init_db.sql")
}

func Dev() error {
	fmt.Println("开发模式运行...")
	return sh.RunV("go", "run", mainPath)
}

func getVersion() string {
	version, _ := sh.Output("git", "describe", "--tags", "--always", "--dirty")
	if version == "" {
		version = "dev"
	}
	return version
}

func getCommit() string {
	commit, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	if commit == "" {
		commit = "none"
	}
	return commit
}
