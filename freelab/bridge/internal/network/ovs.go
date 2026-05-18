package network

import (
	"fmt"
	"os/exec"
)

// CreateOVSBridge creates a new Open vSwitch bridge.
// This represents a physical cable connection or an L2 segment in our simulation.
func CreateOVSBridge(bridgeName string) error {
	cmd := exec.Command("ovs-vsctl", "--may-exist", "add-br", bridgeName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create OVS bridge %s: %v", bridgeName, err)
	}

	// Bring the bridge interface UP via netlink `ip link set dev <br> up`
	cmdUp := exec.Command("ip", "link", "set", "dev", bridgeName, "up")
	if err := cmdUp.Run(); err != nil {
		return fmt.Errorf("failed to bring OVS bridge %s up: %v", bridgeName, err)
	}

	return nil
}

// DeleteOVSBridge removes the Open vSwitch bridge.
// This simulates unplugging a cable completely.
func DeleteOVSBridge(bridgeName string) error {
	cmd := exec.Command("ovs-vsctl", "--if-exists", "del-br", bridgeName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete OVS bridge %s: %v", bridgeName, err)
	}
	return nil
}

// AddPortToBridge connects a TAP interface (Device Port) to an OVS Bridge (Cable).
func AddPortToBridge(bridgeName, tapName string) error {
	cmd := exec.Command("ovs-vsctl", "--may-exist", "add-port", bridgeName, tapName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add port %s to bridge %s: %v", tapName, bridgeName, err)
	}
	return nil
}
