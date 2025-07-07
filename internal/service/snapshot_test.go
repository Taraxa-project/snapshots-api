package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/taraxa/snapshots-api/internal/models"
)

func TestSnapshotService_processSnapshots(t *testing.T) {
	service := NewSnapshotService("test-bucket", "https://test.example.com")

	// Create test snapshots
	snapshots := []*models.Snapshot{
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     100,
			Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-100.tar.gz",
			Filename:  "mainnet-full-100.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     200,
			Timestamp: time.Date(2025, 7, 6, 11, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-200.tar.gz",
			Filename:  "mainnet-full-200.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeLight,
			Block:     150,
			Timestamp: time.Date(2025, 7, 6, 10, 30, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-light-150.tar.gz",
			Filename:  "mainnet-light-150.tar.gz",
		},
		{
			Network:   models.NetworkTestnet,
			Type:      models.SnapshotTypeFull,
			Block:     50,
			Timestamp: time.Date(2025, 7, 6, 9, 0, 0, 0, time.UTC),
			URL:       "https://example.com/testnet-full-50.tar.gz",
			Filename:  "testnet-full-50.tar.gz",
		},
	}

	result := service.processSnapshots(snapshots)

	// Check mainnet results
	mainnetResult, exists := result[models.NetworkMainnet]
	if !exists {
		t.Error("Expected mainnet snapshots to exist")
	}

	if mainnetResult.Full == nil {
		t.Error("Expected mainnet full snapshot to exist")
	} else if mainnetResult.Full.Block != 200 {
		t.Errorf("Expected mainnet full snapshot block 200, got %d", mainnetResult.Full.Block)
	}

	if mainnetResult.Light == nil {
		t.Error("Expected mainnet light snapshot to exist")
	} else if mainnetResult.Light.Block != 150 {
		t.Errorf("Expected mainnet light snapshot block 150, got %d", mainnetResult.Light.Block)
	}

	// Check testnet results
	testnetResult, exists := result[models.NetworkTestnet]
	if !exists {
		t.Error("Expected testnet snapshots to exist")
	}

	if testnetResult.Full == nil {
		t.Error("Expected testnet full snapshot to exist")
	} else if testnetResult.Full.Block != 50 {
		t.Errorf("Expected testnet full snapshot block 50, got %d", testnetResult.Full.Block)
	}

	if testnetResult.Light != nil {
		t.Error("Expected testnet light snapshot to be nil")
	}
}

func TestSnapshotService_findLatestSnapshot(t *testing.T) {
	service := NewSnapshotService("test-bucket", "https://test.example.com")

	tests := []struct {
		name          string
		snapshots     []*models.Snapshot
		expectedNil   bool
		expectedBlock int64
	}{
		{
			name:        "empty slice",
			snapshots:   []*models.Snapshot{},
			expectedNil: true,
		},
		{
			name: "single snapshot",
			snapshots: []*models.Snapshot{
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
				},
			},
			expectedNil:   false,
			expectedBlock: 100,
		},
		{
			name: "multiple snapshots - highest block wins",
			snapshots: []*models.Snapshot{
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC),
				},
				{
					Block:     200,
					Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
				},
				{
					Block:     150,
					Timestamp: time.Date(2025, 7, 6, 11, 0, 0, 0, time.UTC),
				},
			},
			expectedNil:   false,
			expectedBlock: 200,
		},
		{
			name: "same block - latest timestamp wins",
			snapshots: []*models.Snapshot{
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
				},
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC),
				},
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 11, 0, 0, 0, time.UTC),
				},
			},
			expectedNil:   false,
			expectedBlock: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.findLatestSnapshot(tt.snapshots)

			if tt.expectedNil {
				if result != nil {
					t.Error("Expected nil result")
				}
				return
			}

			if result == nil {
				t.Error("Expected non-nil result")
				return
			}

			if result.Block != tt.expectedBlock {
				t.Errorf("Expected block %d, got %d", tt.expectedBlock, result.Block)
			}

			// For same block tests, verify the latest timestamp is selected
			if tt.name == "same block - latest timestamp wins" {
				expectedTime := time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC)
				if !result.Timestamp.Equal(expectedTime) {
					t.Errorf("Expected timestamp %v, got %v", expectedTime, result.Timestamp)
				}
			}
		})
	}
}

func TestSnapshotService_IsValidNetwork(t *testing.T) {
	service := NewSnapshotService("test-bucket", "https://test.example.com")

	tests := []struct {
		network string
		valid   bool
	}{
		{"mainnet", true},
		{"testnet", true},
		{"devnet", true},
		{"invalid", false},
		{"", false},
		{"MAINNET", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.network, func(t *testing.T) {
			result := service.IsValidNetwork(tt.network)
			if result != tt.valid {
				t.Errorf("IsValidNetwork(%s) = %v, want %v", tt.network, result, tt.valid)
			}
		})
	}
}

func TestSnapshotService_GetAllNetworks(t *testing.T) {
	service := NewSnapshotService("test-bucket", "https://test.example.com")

	networks := service.GetAllNetworks()

	expected := []models.Network{
		models.NetworkMainnet,
		models.NetworkTestnet,
		models.NetworkDevnet,
	}

	if len(networks) != len(expected) {
		t.Errorf("Expected %d networks, got %d", len(expected), len(networks))
		return
	}

	for i, network := range networks {
		if network != expected[i] {
			t.Errorf("Expected network %s at index %d, got %s", expected[i], i, network)
		}
	}
}

func TestSnapshotService_fetchSnapshots_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := NewSnapshotService("test-bucket", server.URL)

	_, err := service.fetchSnapshots()
	if err == nil {
		t.Error("Expected error from fetchSnapshots")
	}
}

func TestSnapshotService_fetchSnapshots_Success(t *testing.T) {
	// Create a test server that returns valid GCP response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{
			"kind": "storage#objects",
			"items": [
				{"name": "mainnet-full-db-block-19547931-20250706-062734.tar.gz"},
				{"name": "testnet-light-db-block-2516167-20250706-052226.tar.gz"},
				{"name": "invalid-file.txt"}
			]
		}`
		w.Write([]byte(response))
	}))
	defer server.Close()

	service := NewSnapshotService("test-bucket", server.URL)

	snapshots, err := service.fetchSnapshots()
	if err != nil {
		t.Errorf("Unexpected error from fetchSnapshots: %v", err)
		return
	}

	// Should have 2 valid snapshots (invalid-file.txt should be skipped)
	if len(snapshots) != 2 {
		t.Errorf("Expected 2 snapshots, got %d", len(snapshots))
	}

	// Verify first snapshot
	if snapshots[0].Network != models.NetworkMainnet {
		t.Errorf("Expected mainnet, got %s", snapshots[0].Network)
	}
	if snapshots[0].Type != models.SnapshotTypeFull {
		t.Errorf("Expected full, got %s", snapshots[0].Type)
	}
	if snapshots[0].Block != 19547931 {
		t.Errorf("Expected block 19547931, got %d", snapshots[0].Block)
	}
}
