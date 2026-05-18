package state

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"freelab-bridge/internal/models"
)

// TopologyManager handles thread-safe operations on network topologies in memory.
type TopologyManager struct {
	mu         sync.RWMutex
	topologies map[string]*models.Topology
}

// NewManager initializes a new TopologyManager.
func NewManager() *TopologyManager {
	return &TopologyManager{
		topologies: make(map[string]*models.Topology),
	}
}

// CreateTopology initializes a new empty topology and stores it.
func (m *TopologyManager) CreateTopology(name, mapName string) *models.Topology {
	m.mu.Lock()
	defer m.mu.Unlock()

	topo := &models.Topology{
		Version:  "1.0",
		ID:       "topo_" + uuid.New().String()[:8], // Short UUID for readability
		Name:     name,
		Map:      mapName,
		Created:  time.Now(),
		Modified: time.Now(),
		Devices:  []models.Device{},
		Links:    []models.Link{},
	}

	m.topologies[topo.ID] = topo
	return topo
}

// GetTopology retrieves a topology by its ID.
func (m *TopologyManager) GetTopology(id string) (*models.Topology, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	topo, exists := m.topologies[id]
	if !exists {
		return nil, errors.New("topology not found")
	}
	return topo, nil
}

// AddDevice adds a new device to an existing topology.
func (m *TopologyManager) AddDevice(topoID string, device models.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	topo, exists := m.topologies[topoID]
	if !exists {
		return errors.New("topology not found")
	}

	device.State = "defined" // Initial state before VM is booted
	topo.Devices = append(topo.Devices, device)
	topo.Modified = time.Now()
	return nil
}

// AddLink adds a new physical cable connection between two interfaces.
func (m *TopologyManager) AddLink(topoID string, link models.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	topo, exists := m.topologies[topoID]
	if !exists {
		return errors.New("topology not found")
	}

	link.ID = "link_" + uuid.New().String()[:8]
	link.State = "defined"
	topo.Links = append(topo.Links, link)
	topo.Modified = time.Now()
	return nil
}
