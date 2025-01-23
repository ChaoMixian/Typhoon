// File: daemon.go
package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"Typhoon/config"
	"Typhoon/utils"
)

// 全局变量
var (
	mihomoProcess *exec.Cmd
	outputReader  *io.PipeReader
	outputWriter  *io.PipeWriter
	mu            sync.Mutex
)

// StartMihomo ...
func StartMihomo(binaryPath string, args []string) error {
	mu.Lock()
	defer mu.Unlock()

	// 判断是否已经有进程在跑
	if mihomoProcess != nil && mihomoProcess.Process != nil {
		return fmt.Errorf("mihomo is already running")
	}

	// 创建管道，用来获取子进程输出
	outputReader, outputWriter = io.Pipe()
	mihomoProcess = exec.Command(binaryPath, args...)
	mihomoProcess.Stdout = outputWriter
	mihomoProcess.Stderr = outputWriter

	// 启动进程
	if err := mihomoProcess.Start(); err != nil {
		return fmt.Errorf("failed to start Mihomo: %v", err)
	}

	// 异步读取输出
	go captureOutput()

	log.Printf("Mihomo started with PID: %d", mihomoProcess.Process.Pid)
	return nil
}

// StopMihomo ...
func StopMihomo() error {
	mu.Lock()
	defer mu.Unlock()

	if mihomoProcess == nil || mihomoProcess.Process == nil {
		return fmt.Errorf("mihomo is not running")
	}

	// 发 SIGTERM 以优雅停止
	if err := mihomoProcess.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop Mihomo: %v", err)
	}

	// 等待进程退出
	if err := mihomoProcess.Wait(); err != nil {
		mihomoProcess = nil
		return fmt.Errorf("failed to wait for Mihomo termination: %v", err)
	}

	log.Println("Mihomo stopped successfully")
	mihomoProcess = nil
	return nil
}

// RestartMihomo ...
func RestartMihomo(binaryPath string, args []string) error {
	if err := StopMihomo(); err != nil {
		return fmt.Errorf("failed to stop Mihomo: %v", err)
	}
	return StartMihomo(binaryPath, args)
}

// captureOutput ...
func captureOutput() {
	for {
		buf := make([]byte, 1024)
		n, err := outputReader.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Error reading Mihomo output: %v", err)
			break
		}
		if n > 0 {
			log.Printf("Mihomo output: %s", string(buf[:n]))
		}
	}
}

// ==================== 基于配置启动/重启的函数 ====================

// StartMihomoFromConfig 加载配置、校验文件、Patch 配置并启动
func StartMihomoFromConfig() error {
	// 1. 加载全局配置
	cfg := config.GetConfig(config.ConfigFilePath, false)

	binaryPath := cfg.Proxy.Mihomo.BinPath
	// configPath := cfg.Proxy.Mihomo.ConfigPath
	// runtimeConfigPath := cfg.Proxy.Mihomo.RuntimeConfigPath
	configPath, runtimeConfigPath := utils.GetMihomoConfigPath(cfg.Proxy.Mihomo.CurrentConfig)
	log.Printf("configPath: %s, runtimeConfigPath: %s", configPath, runtimeConfigPath)
	controllerAddress := cfg.Proxy.Mihomo.ControllerAddress
	port := cfg.Proxy.Mihomo.ListenPort

	// 2. 校验文件是否存在
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary file not found: %v", err)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %v", err)
	}

	// 3. Patch 配置（先修改 configPath -> runtimeConfigPath）
	if err := config.PatchMihomoConfig(configPath, runtimeConfigPath, controllerAddress, port); err != nil {
		return fmt.Errorf("failed to patch Mihomo config: %v", err)
	}

	// 4. 组装启动参数
	args := []string{"-f", runtimeConfigPath}

	// 5. 调用原有的 StartMihomo
	if err := StartMihomo(binaryPath, args); err != nil {
		return err
	}
	return nil
}

// RestartMihomoFromConfig 重启，等同于 Stop + StartMihomoFromConfig
func RestartMihomoFromConfig() error {
	// 先停止
	if err := StopMihomo(); err != nil {
		return fmt.Errorf("failed to stop Mihomo: %v", err)
	}
	// 再启动
	return StartMihomoFromConfig()
}

func IsMihomoRunning() bool {
	mu.Lock()
	defer mu.Unlock()

	return mihomoProcess != nil && mihomoProcess.Process != nil
}
