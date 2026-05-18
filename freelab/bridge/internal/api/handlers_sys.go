package api

import (
	"encoding/json"
	"net/http"
	"os"
)

// SystemStatusResponse maps to the API specification for system health.
type SystemStatusResponse struct {
	KVMAvailable      bool `json:"kvm_available"`
	TotalRAMMb        int  `json:"total_ram_mb"`
	AvailableRAMMb    int  `json:"available_ram_mb"`
	CPUCores          int  `json:"cpu_cores"`
	RunningVMs        int  `json:"running_vms"`
	RunningContainers int  `json:"running_containers"`
}

// HandleSystemStatus returns the current capacity and state of the hypervisor host.
func HandleSystemStatus(w http.ResponseWriter, r *http.Request) {
	// Simple check to see if KVM is exposed to WSL2
	kvmEnabled := false
	if _, err := os.Stat("/dev/kvm"); err == nil {
		kvmEnabled = true
	}

	// TODO: Replace memory and CPU static values with real OS metrics 
	// using a library like 'shirou/gopsutil/v3' later.
	response := SystemStatusResponse{
		KVMAvailable:      kvmEnabled,
		TotalRAMMb:        16384, // Mocked 16GB
		AvailableRAMMb:    8192,  // Mocked 8GB
		CPUCores:          8,     // Mocked 8 cores
		RunningVMs:        0,
		RunningContainers: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
