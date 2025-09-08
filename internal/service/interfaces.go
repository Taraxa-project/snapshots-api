package service

import "github.com/taraxa/snapshots-api/internal/models"

// SnapshotServiceInterface defines the contract for snapshot service
type SnapshotServiceInterface interface {
	GetSnapshots(network models.Network) (*models.NetworkSnapshots, error)
	GetSnapshotsWithAuth(network models.Network, authenticated bool) (*models.NetworkSnapshots, error)
	IsValidNetwork(network string) bool
	GetAllNetworks() []models.Network
}
