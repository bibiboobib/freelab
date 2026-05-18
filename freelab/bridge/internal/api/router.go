package api

import (
	"github.com/go-chi/chi/v5"
	"freelab-bridge/internal/state"
	"freelab-bridge/internal/qemu" // ПРИВЯЗЫВАЕМ НАШ QEMU
)

// NewRouter initializes and configures the main HTTP router.
func NewRouter(stateManager *state.TopologyManager, qemuManager *qemu.QemuManager) *chi.Mux {
	r := chi.NewRouter()

	// Initialize API Handlers with dependencies
	apiHandler := &Handler{
		StateManager: stateManager,
		QemuManager:  qemuManager, // ПЕРЕДАЕМ ЗАВИСИМОСТЬ
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/system/status", HandleSystemStatus)

		r.Route("/topology", func(r chi.Router) {
			r.Post("/", apiHandler.HandleCreateTopology)
			r.Get("/{id}", apiHandler.HandleGetTopology)
			r.Post("/{id}/devices", apiHandler.HandleAddDevice)
			// r.Post("/{id}/links", apiHandler.HandleAddLink) // We will wire this later
		})
	})

	return r
}
