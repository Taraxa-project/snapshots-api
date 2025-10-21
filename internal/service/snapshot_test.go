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

	// Create test snapshots with multiple snapshots per type to test previous arrays
	snapshots := []*models.Snapshot{
		// Mainnet full snapshots (5 total)
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     500,
			Timestamp: time.Date(2025, 7, 10, 10, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-500.tar.gz",
			Filename:  "mainnet-full-500.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     400,
			Timestamp: time.Date(2025, 7, 9, 10, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-400.tar.gz",
			Filename:  "mainnet-full-400.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     300,
			Timestamp: time.Date(2025, 7, 8, 10, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-300.tar.gz",
			Filename:  "mainnet-full-300.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     200,
			Timestamp: time.Date(2025, 7, 7, 10, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-200.tar.gz",
			Filename:  "mainnet-full-200.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeFull,
			Block:     100,
			Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-full-100.tar.gz",
			Filename:  "mainnet-full-100.tar.gz",
		},
		// Mainnet light snapshots (4 total)
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeLight,
			Block:     450,
			Timestamp: time.Date(2025, 7, 9, 15, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-light-450.tar.gz",
			Filename:  "mainnet-light-450.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeLight,
			Block:     350,
			Timestamp: time.Date(2025, 7, 8, 15, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-light-350.tar.gz",
			Filename:  "mainnet-light-350.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeLight,
			Block:     250,
			Timestamp: time.Date(2025, 7, 7, 15, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-light-250.tar.gz",
			Filename:  "mainnet-light-250.tar.gz",
		},
		{
			Network:   models.NetworkMainnet,
			Type:      models.SnapshotTypeLight,
			Block:     150,
			Timestamp: time.Date(2025, 7, 6, 15, 0, 0, 0, time.UTC),
			URL:       "https://example.com/mainnet-light-150.tar.gz",
			Filename:  "mainnet-light-150.tar.gz",
		},
		// Testnet full snapshot (only 1)
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

	// Check mainnet full snapshots
	if mainnetResult.Full == nil {
		t.Error("Expected mainnet full snapshot to exist")
	} else if mainnetResult.Full.Block != 500 {
		t.Errorf("Expected mainnet full snapshot block 500, got %d", mainnetResult.Full.Block)
	}

	// Check mainnet previous full snapshots
	if len(mainnetResult.PreviousFull) != 3 {
		t.Errorf("Expected 3 previous full snapshots, got %d", len(mainnetResult.PreviousFull))
	} else {
		expectedBlocks := []int64{400, 300, 200}
		for i, expected := range expectedBlocks {
			if mainnetResult.PreviousFull[i].Block != expected {
				t.Errorf("Expected previous full snapshot %d to have block %d, got %d", i, expected, mainnetResult.PreviousFull[i].Block)
			}
		}
	}

	// Check mainnet light snapshots
	if mainnetResult.Light == nil {
		t.Error("Expected mainnet light snapshot to exist")
	} else if mainnetResult.Light.Block != 450 {
		t.Errorf("Expected mainnet light snapshot block 450, got %d", mainnetResult.Light.Block)
	}

	// Check mainnet previous light snapshots
	if len(mainnetResult.PreviousLight) != 3 {
		t.Errorf("Expected 3 previous light snapshots, got %d", len(mainnetResult.PreviousLight))
	} else {
		expectedBlocks := []int64{350, 250, 150}
		for i, expected := range expectedBlocks {
			if mainnetResult.PreviousLight[i].Block != expected {
				t.Errorf("Expected previous light snapshot %d to have block %d, got %d", i, expected, mainnetResult.PreviousLight[i].Block)
			}
		}
	}

	// Check testnet results
	testnetResult, exists := result[models.NetworkTestnet]
	if !exists {
		t.Error("Expected testnet snapshots to exist")
	}

	// Testnet should have only 1 full snapshot (no previous)
	if testnetResult.Full == nil {
		t.Error("Expected testnet full snapshot to exist")
	} else if testnetResult.Full.Block != 50 {
		t.Errorf("Expected testnet full snapshot block 50, got %d", testnetResult.Full.Block)
	}

	// Testnet should have no previous full snapshots
	if len(testnetResult.PreviousFull) != 0 {
		t.Errorf("Expected 0 previous full snapshots for testnet, got %d", len(testnetResult.PreviousFull))
	}

	// Testnet should have no light snapshots
	if testnetResult.Light != nil {
		t.Error("Expected testnet light snapshot to be nil")
	}

	if len(testnetResult.PreviousLight) != 0 {
		t.Errorf("Expected 0 previous light snapshots for testnet, got %d", len(testnetResult.PreviousLight))
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

func TestSnapshotService_findLatestAndPreviousSnapshots(t *testing.T) {
	service := NewSnapshotService("test-bucket", "https://test.example.com")

	tests := []struct {
		name                  string
		snapshots             []*models.Snapshot
		expectedLatestNil     bool
		expectedLatestBlock   int64
		expectedPreviousCount int
		expectedPreviousBlocks []int64
	}{
		{
			name:                  "empty slice",
			snapshots:             []*models.Snapshot{},
			expectedLatestNil:     true,
			expectedPreviousCount: 0,
		},
		{
			name: "single snapshot",
			snapshots: []*models.Snapshot{
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
				},
			},
			expectedLatestNil:     false,
			expectedLatestBlock:   100,
			expectedPreviousCount: 0,
		},
		{
			name: "two snapshots",
			snapshots: []*models.Snapshot{
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
				},
				{
					Block:     200,
					Timestamp: time.Date(2025, 7, 6, 11, 0, 0, 0, time.UTC),
				},
			},
			expectedLatestNil:      false,
			expectedLatestBlock:    200,
			expectedPreviousCount:  1,
			expectedPreviousBlocks: []int64{100},
		},
		{
			name: "five snapshots - should return latest + 3 previous",
			snapshots: []*models.Snapshot{
				{
					Block:     100,
					Timestamp: time.Date(2025, 7, 6, 10, 0, 0, 0, time.UTC),
				},
				{
					Block:     200,
					Timestamp: time.Date(2025, 7, 6, 11, 0, 0, 0, time.UTC),
				},
				{
					Block:     300,
					Timestamp: time.Date(2025, 7, 6, 12, 0, 0, 0, time.UTC),
				},
				{
					Block:     400,
					Timestamp: time.Date(2025, 7, 6, 13, 0, 0, 0, time.UTC),
				},
				{
					Block:     500,
					Timestamp: time.Date(2025, 7, 6, 14, 0, 0, 0, time.UTC),
				},
			},
			expectedLatestNil:      false,
			expectedLatestBlock:    500,
			expectedPreviousCount:  3,
			expectedPreviousBlocks: []int64{400, 300, 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			latest, previous := service.findLatestAndPreviousSnapshots(tt.snapshots)

			// Check latest
			if tt.expectedLatestNil {
				if latest != nil {
					t.Error("Expected nil latest result")
				}
			} else {
				if latest == nil {
					t.Error("Expected non-nil latest result")
					return
				}
				if latest.Block != tt.expectedLatestBlock {
					t.Errorf("Expected latest block %d, got %d", tt.expectedLatestBlock, latest.Block)
				}
			}

			// Check previous count
			if len(previous) != tt.expectedPreviousCount {
				t.Errorf("Expected %d previous snapshots, got %d", tt.expectedPreviousCount, len(previous))
			}

			// Check previous blocks
			if tt.expectedPreviousBlocks != nil {
				for i, expectedBlock := range tt.expectedPreviousBlocks {
					if i >= len(previous) {
						t.Errorf("Expected previous snapshot at index %d with block %d, but not enough snapshots", i, expectedBlock)
						continue
					}
					if previous[i].Block != expectedBlock {
						t.Errorf("Expected previous snapshot %d to have block %d, got %d", i, expectedBlock, previous[i].Block)
					}
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
