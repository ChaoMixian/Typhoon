// File: core.go
// manage the mihomo core or xray core, including start, stop, restart, etc.

package core

import (
	"fmt"
	"os/exec"
)

// MihomoManager 管理 Mihomo 内核的结构体
type MihomoManager struct {
	// 可根据需求扩展其他字段，比如配置等
}

// NewMihomoManager 创建一个新的 Mihomo 管理器实例
func NewMihomoManager() *MihomoManager {
	return &MihomoManager{}
}

// Start 启动 Mihomo 内核
func (m *MihomoManager) Start() error {
	// 假设mihomo是通过命令启动的
	cmd := exec.Command("mihomo", "--start")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start mihomo: %v", err)
	}
	return nil
}

// Stop 停止 Mihomo 内核
func (m *MihomoManager) Stop() error {
	cmd := exec.Command("mihomo", "--stop")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to stop mihomo: %v", err)
	}
	return nil
}

// Status 获取 Mihomo 内核的状态
func (m *MihomoManager) Status() (string, error) {
	cmd := exec.Command("mihomo", "--status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get mihomo status: %v", err)
	}
	return string(output), nil
}
