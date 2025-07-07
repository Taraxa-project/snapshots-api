package api

import (
	"github.com/taraxa/snapshots-api/internal/models"
)

// MockSnapshotService is a mock implementation for testing
type MockSnapshotService struct {
	GetSnapshotsFunc   func(network models.Network) (*models.NetworkSnapshots, error)
	IsValidNetworkFunc func(network string) bool
	GetAllNetworksFunc func() []models.Network
}

func (m *MockSnapshotService) GetSnapshots(network models.Network) (*models.NetworkSnapshots, error) {
	if m.GetSnapshotsFunc != nil {
		return m.GetSnapshotsFunc(network)
	}
	// Default implementation
	return &models.NetworkSnapshots{
		Full: &models.SnapshotInfo{
			Block:     12345,
			Timestamp: "2025-07-06 14:30",
			URL:       "https://storage.googleapis.com/taraxa-snapshot/mainnet-full-db-block-12345-20250706-143000.tar.gz",
		},
		Light: &models.SnapshotInfo{
			Block:     12345,
			Timestamp: "2025-07-06 14:30",
			URL:       "https://storage.googleapis.com/taraxa-snapshot/mainnet-light-db-block-12345-20250706-143000.tar.gz",
		},
	}, nil
}

func (m *MockSnapshotService) IsValidNetwork(network string) bool {
	if m.IsValidNetworkFunc != nil {
		return m.IsValidNetworkFunc(network)
	}
	// Default implementation
	switch network {
	case "mainnet", "testnet", "devnet":
		return true
	default:
		return false
	}
}

func (m *MockSnapshotService) GetAllNetworks() []models.Network {
	if m.GetAllNetworksFunc != nil {
		return m.GetAllNetworksFunc()
	}
	// Default implementation
	return []models.Network{
		models.NetworkMainnet,
		models.NetworkTestnet,
		models.NetworkDevnet,
	}
}
