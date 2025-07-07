package models

import (
	"testing"
	"time"
)

func TestSnapshot_ToSnapshotInfo(t *testing.T) {
	timestamp := time.Date(2025, 7, 6, 14, 30, 45, 0, time.UTC)

	snapshot := &Snapshot{
		Network:   NetworkMainnet,
		Type:      SnapshotTypeFull,
		Block:     12345,
		Timestamp: timestamp,
		URL:       "https://example.com/snapshot.tar.gz",
		Filename:  "snapshot.tar.gz",
	}

	result := snapshot.ToSnapshotInfo()

	if result.Block != 12345 {
		t.Errorf("Expected block 12345, got %d", result.Block)
	}

	expectedTimestamp := "2025-07-06 14:30"
	if result.Timestamp != expectedTimestamp {
		t.Errorf("Expected timestamp %s, got %s", expectedTimestamp, result.Timestamp)
	}

	if result.URL != "https://example.com/snapshot.tar.gz" {
		t.Errorf("Expected URL %s, got %s", "https://example.com/snapshot.tar.gz", result.URL)
	}
}

func TestNetworkConstants(t *testing.T) {
	tests := []struct {
		network  Network
		expected string
	}{
		{NetworkMainnet, "mainnet"},
		{NetworkTestnet, "testnet"},
		{NetworkDevnet, "devnet"},
	}

	for _, tt := range tests {
		t.Run(string(tt.network), func(t *testing.T) {
			if string(tt.network) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.network))
			}
		})
	}
}

func TestSnapshotTypeConstants(t *testing.T) {
	tests := []struct {
		snapshotType SnapshotType
		expected     string
	}{
		{SnapshotTypeFull, "full"},
		{SnapshotTypeLight, "light"},
	}

	for _, tt := range tests {
		t.Run(string(tt.snapshotType), func(t *testing.T) {
			if string(tt.snapshotType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.snapshotType))
			}
		})
	}
}
