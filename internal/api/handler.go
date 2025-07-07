package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/taraxa/snapshots-api/internal/models"
	"github.com/taraxa/snapshots-api/internal/service"
)

// Handler holds the API handlers
type Handler struct {
	snapshotService service.SnapshotServiceInterface
}

// NewHandler creates a new API handler
func NewHandler(snapshotService service.SnapshotServiceInterface) *Handler {
	return &Handler{
		snapshotService: snapshotService,
	}
}

// Routes sets up the HTTP routes
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.getSnapshots)
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/ready", h.ready)

	return mux
}

// getSnapshots handles GET requests for snapshot data
func (h *Handler) getSnapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get network parameter
	network := r.URL.Query().Get("network")
	if network == "" {
		http.Error(w, "network parameter is required", http.StatusBadRequest)
		return
	}

	// Validate network
	if !h.snapshotService.IsValidNetwork(network) {
		http.Error(w, "invalid network. Supported networks: mainnet, testnet, devnet", http.StatusBadRequest)
		return
	}

	// Get snapshots
	snapshots, err := h.snapshotService.GetSnapshots(models.Network(network))
	if err != nil {
		log.Printf("Error fetching snapshots for network %s: %v", network, err)
		http.Error(w, "failed to fetch snapshots", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes

	// Encode and send response
	if err := json.NewEncoder(w).Encode(snapshots); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// health handles health check requests
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "healthy",
		"service": "snapshots-api",
	}

	json.NewEncoder(w).Encode(response)
}

// ready handles readiness check requests
func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Try to fetch snapshots to verify service is ready
	_, err := h.snapshotService.GetSnapshots(models.NetworkMainnet)
	if err != nil {
		log.Printf("Readiness check failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)

		response := map[string]string{
			"status": "not ready",
			"error":  "failed to connect to GCP bucket",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ready",
		"service": "snapshots-api",
	}

	json.NewEncoder(w).Encode(response)
}
