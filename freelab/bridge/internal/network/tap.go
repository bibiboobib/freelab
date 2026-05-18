package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

// CreateTap creates a new Layer 2 TAP interface for QEMU to bind to.
// Typical name format: nl_R1_e0 (freelab_Router1_eth0)
func CreateTap(name string) error {
	la := netlink.NewLinkAttrs()
	la.Name = name

	// Tuntap interface configured as TAP (Layer 2)
	tuntap := &netlink.Tuntap{
		LinkAttrs: la,
		Mode:      netlink.TUNTAP_MODE_TAP,
		Queues:    1,
		Flags:     netlink.TUNTAP_ONE_QUEUE | netlink.TUNTAP_VNET_HDR,
	}

	if err := netlink.LinkAdd(tuntap); err != nil {
		return fmt.Errorf("failed to create TAP %s: %v", name, err)
	}

	// Bring the interface UP (equivalent to 'ip link set dev name up')
	if err := netlink.LinkSetUp(tuntap); err != nil {
		return fmt.Errorf("failed to bring TAP %s up: %v", name, err)
	}

	return nil
}

// DeleteTap completely removes the TAP interface from the Linux kernel.
func DeleteTap(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("TAP %s not found: %v", name, err)
	}

	if err := netlink.LinkDel(link); err != nil {
		return fmt.Errorf("failed to delete TAP %s: %v", name, err)
	}

	return nil
}
