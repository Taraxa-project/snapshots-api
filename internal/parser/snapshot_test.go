package parser

import (
	"testing"
	"time"

	"github.com/taraxa/snapshots-api/internal/models"
)

func TestSnapshotParser_ParseSnapshot(t *testing.T) {
	parser := NewSnapshotParser()
	baseURL := "https://storage.googleapis.com/taraxa-snapshot"

	tests := []struct {
		name     string
		filename string
		expected *models.Snapshot
		wantErr  bool
	}{
		{
			name:     "valid mainnet full snapshot",
			filename: "mainnet-full-db-block-19547931-20250706-062734.tar.gz",
			expected: &models.Snapshot{
				Network:   models.NetworkMainnet,
				Type:      models.SnapshotTypeFull,
				Block:     19547931,
				Timestamp: time.Date(2025, 7, 6, 6, 27, 34, 0, time.UTC),
				URL:       "https://storage.googleapis.com/taraxa-snapshot/mainnet-full-db-block-19547931-20250706-062734.tar.gz",
				Filename:  "mainnet-full-db-block-19547931-20250706-062734.tar.gz",
			},
			wantErr: false,
		},
		{
			name:     "valid testnet light snapshot",
			filename: "testnet-light-db-block-2516167-20250706-052226.tar.gz",
			expected: &models.Snapshot{
				Network:   models.NetworkTestnet,
				Type:      models.SnapshotTypeLight,
				Block:     2516167,
				Timestamp: time.Date(2025, 7, 6, 5, 22, 26, 0, time.UTC),
				URL:       "https://storage.googleapis.com/taraxa-snapshot/testnet-light-db-block-2516167-20250706-052226.tar.gz",
				Filename:  "testnet-light-db-block-2516167-20250706-052226.tar.gz",
			},
			wantErr: false,
		},
		{
			name:     "valid devnet full snapshot",
			filename: "devnet-full-db-block-394662-20250706-042052.tar.gz",
			expected: &models.Snapshot{
				Network:   models.NetworkDevnet,
				Type:      models.SnapshotTypeFull,
				Block:     394662,
				Timestamp: time.Date(2025, 7, 6, 4, 20, 52, 0, time.UTC),
				URL:       "https://storage.googleapis.com/taraxa-snapshot/devnet-full-db-block-394662-20250706-042052.tar.gz",
				Filename:  "devnet-full-db-block-394662-20250706-042052.tar.gz",
			},
			wantErr: false,
		},
		{
			name:     "invalid filename format",
			filename: "invalid-filename.tar.gz",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid network",
			filename: "invalid-full-db-block-123-20250706-042052.tar.gz",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid snapshot type",
			filename: "mainnet-invalid-db-block-123-20250706-042052.tar.gz",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid block number",
			filename: "mainnet-full-db-block-abc-20250706-042052.tar.gz",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid timestamp",
			filename: "mainnet-full-db-block-123-20250706-999999.tar.gz",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseSnapshot(tt.filename, baseURL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseSnapshot() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSnapshot() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("ParseSnapshot() returned nil result")
				return
			}

			// Compare fields
			if result.Network != tt.expected.Network {
				t.Errorf("Network = %v, want %v", result.Network, tt.expected.Network)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}
			if result.Block != tt.expected.Block {
				t.Errorf("Block = %v, want %v", result.Block, tt.expected.Block)
			}
			if !result.Timestamp.Equal(tt.expected.Timestamp) {
				t.Errorf("Timestamp = %v, want %v", result.Timestamp, tt.expected.Timestamp)
			}
			if result.URL != tt.expected.URL {
				t.Errorf("URL = %v, want %v", result.URL, tt.expected.URL)
			}
			if result.Filename != tt.expected.Filename {
				t.Errorf("Filename = %v, want %v", result.Filename, tt.expected.Filename)
			}
		})
	}
}

func TestSnapshotParser_IsValidNetwork(t *testing.T) {
	parser := NewSnapshotParser()

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
			result := parser.IsValidNetwork(tt.network)
			if result != tt.valid {
				t.Errorf("IsValidNetwork(%s) = %v, want %v", tt.network, result, tt.valid)
			}
		})
	}
}
