// modify iptables rules
package daemon

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// EnableIPForwarding enables IP forwarding on the system
func EnableIPForwarding() error {
	// Check if IP forwarding is already enabled
	cmd := exec.Command("sysctl", "-n", "net.ipv4.ip_forward")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check IP forwarding status: %v", err)
	}
	if strings.TrimSpace(string(output)) != "1" {
		// Enable IP forwarding
		cmd = exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1")
		output, err = cmd.CombinedOutput()
		log.Println(string(output))
		if err != nil {
			return fmt.Errorf("failed to enable IP forwarding: %v", err)
		}
		log.Println("IP forwarding enabled successfully.")
	} else {
		log.Println("IP forwarding is already enabled.")
	}
	return nil
}

// ConfigureIptables configures the iptables for transparent proxying
// listenPort: the port where the proxy is listening (e.g., 1080)
// proxyPort: the port to which traffic will be forwarded (e.g., 7890)
func ConfigureIptables(listenPort, proxyPort int) error {
	// Flush existing iptables rules to avoid conflicts
	if err := flushIptables(); err != nil {
		return fmt.Errorf("failed to flush iptables: %v", err)
	}

	// Set up the transparent proxy redirection rule
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-p", "tcp", "--dport", fmt.Sprintf("%d", listenPort), "-j", "REDIRECT", "--to-port", fmt.Sprintf("%d", proxyPort))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add iptables rule: %v, output: %s", err, string(output))
	}

	log.Printf("iptables rule added to redirect port %d traffic to port %d\n", listenPort, proxyPort)
	return nil
}

// flushIptables flushes all current iptables rules
func flushIptables() error {
	// Flush all current iptables rules
	cmd := exec.Command("iptables", "-F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to flush iptables: %v, output: %s", err, string(output))
	}
	log.Println("Existing iptables rules flushed successfully.")
	return nil
}

// ClearIptables clears the iptables rules related to transparent proxying
func ClearIptables(listenPort int) error {
	// Remove the transparent proxy rule
	cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-p", "tcp", "--dport", fmt.Sprintf("%d", listenPort), "-j", "REDIRECT", "--to-port", fmt.Sprintf("%d", listenPort))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove iptables rule: %v, output: %s", err, string(output))
	}

	log.Printf("iptables rule for port %d removed successfully\n", listenPort)
	return nil
}
