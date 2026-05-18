package models

import "time"

// Topology represents a complete network laboratory project.
// This directly maps to the JSON structure defined in the specifications.
type Topology struct {
	Version  string    `json:"version"`
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Map      string    `json:"map"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Devices  []Device  `json:"devices"`
	Links    []Link    `json:"links"`
}

// Device describes a single virtual node (VM or Container).
type Device struct {
	ID            string      `json:"id"`
	CatalogID     string      `json:"catalog_id"`
	Hostname      string      `json:"hostname"`
	RackID        string      `json:"rack_id,omitempty"`
	RackPositionU int         `json:"rack_position_u,omitempty"`
	Interfaces    []Interface `json:"interfaces"`
	State         string      `json:"state"` // Expected states: "defined", "running", "stopped", "error"
}

// Interface describes a physical network port on a device.
// In the backend, this correlates to a Linux TAP interface.
type Interface struct {
	Index       int    `json:"index"`
	Name        string `json:"name"` // e.g., "GigabitEthernet0/0"
	IP          string `json:"ip,omitempty"`
	Mask        string `json:"mask,omitempty"`
	VLAN        *int   `json:"vlan,omitempty"`
	Description string `json:"description,omitempty"`
}

// Link represents a physical cable connecting two device interfaces.
// In the backend, this correlates to a Linux Bridge or Open vSwitch.
type Link struct {
	ID        string       `json:"id"`
	CableType string       `json:"cable_type"` // e.g., "ethernet_cat6", "dac"
	EndpointA LinkEndpoint `json:"endpoint_a"`
	EndpointB LinkEndpoint `json:"endpoint_b"`
	State     string       `json:"state"` // Expected states: "up", "down"
}

// LinkEndpoint defines exactly where a cable is plugged in.
type LinkEndpoint struct {
	DeviceID       string `json:"device"`
	InterfaceIndex int    `json:"interface_index"`
}
