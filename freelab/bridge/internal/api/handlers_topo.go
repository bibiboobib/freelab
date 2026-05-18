package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"freelab-bridge/internal/models"
	"freelab-bridge/internal/state"
	"freelab-bridge/internal/qemu" // ДОБАВИЛИ ИМПОРТ
)

// Handler holds dependencies for API routes.
type Handler struct {
	StateManager *state.TopologyManager
	QemuManager  *qemu.QemuManager // ДОБАВИЛИ МЕНЕДЖЕР
}

// ... HandleCreateTopology и HandleGetTopology остаются без изменений ...

type CreateTopologyRequest struct {
	Name string `json:"name"`
	Map  string `json:"map"`
}

func (h *Handler) HandleCreateTopology(w http.ResponseWriter, r *http.Request) {
	var req CreateTopologyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	topo := h.StateManager.CreateTopology(req.Name, req.Map)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topo)
}

func (h *Handler) HandleGetTopology(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	topo, err := h.StateManager.GetTopology(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topo)
}

// HandleAddDevice registers a new 3D device into the logical topology AND STARTS IT!
func (h *Handler) HandleAddDevice(w http.ResponseWriter, r *http.Request) {
	topoID := chi.URLParam(r, "id") // ID топологии из URL
	
	var device models.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, "Invalid device payload", http.StatusBadRequest)
		return
	}

	// 1. ЛОГИКА: Добавляем в оперативную память (State)
	if err := h.StateManager.AddDevice(topoID, device); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 2. ФИЗИКА: Запускаем настоящую виртуалку!
	// Генерируем уникальное имя для гипервизора (Например: "Router1")
	// Если у вас в models.Device нет поля ID, можете использовать device.Name
	vmID := fmt.Sprintf("vm-%s", device.ID) 
	if device.ID == "" { 
		vmID = "vm-unknown" // Запасной вариант
	}

	log.Printf("[API] Начинаем физический запуск устройства %s...", vmID)

	// Делаем мгновенный клон базового диска (пока жестко прописываем cirros.qcow2)
	overlayPath, err := h.QemuManager.CreateOverlay("cirros.qcow2", vmID)
	if err != nil {
		log.Printf("[ОШИБКА] Не удалось создать диск: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create disk: %v", err), http.StatusInternalServerError)
		return
	}

	// Запускаем процесс QEMU
	process, err := h.QemuManager.StartVM(vmID, overlayPath, []string{})
	if err != nil {
		log.Printf("[ОШИБКА] Не удалось запустить QEMU: %v", err)
		http.Error(w, fmt.Sprintf("Failed to start VM: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[API] УСПЕХ! Устройство %s запущено (PID: %d)", vmID, process.Pid)

	w.WriteHeader(http.StatusCreated)
}
