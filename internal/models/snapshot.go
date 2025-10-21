package models

import "time"

// SnapshotType represents the type of snapshot (full or light)
type SnapshotType string

const (
	SnapshotTypeFull  SnapshotType = "full"
	SnapshotTypeLight SnapshotType = "light"
)

// Network represents the blockchain network
type Network string

const (
	NetworkMainnet Network = "mainnet"
	NetworkTestnet Network = "testnet"
	NetworkDevnet  Network = "devnet"
)

// Snapshot represents a single snapshot file
type Snapshot struct {
	Network   Network      `json:"-"`
	Type      SnapshotType `json:"-"`
	Block     int64        `json:"block"`
	Timestamp time.Time    `json:"-"`
	URL       string       `json:"url"`
	Filename  string       `json:"-"`
}

// SnapshotInfo represents the formatted timestamp for API response
type SnapshotInfo struct {
	Block     int64  `json:"block"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
}

// NetworkSnapshots represents snapshots for a specific network
type NetworkSnapshots struct {
	Full          *SnapshotInfo  `json:"full,omitempty"`
	Light         *SnapshotInfo  `json:"light,omitempty"`
	PreviousLight []SnapshotInfo `json:"previous-light,omitempty"`
	PreviousFull  []SnapshotInfo `json:"previous-full,omitempty"`
}

// ToSnapshotInfo converts a Snapshot to SnapshotInfo with formatted timestamp
func (s *Snapshot) ToSnapshotInfo() *SnapshotInfo {
	return &SnapshotInfo{
		Block:     s.Block,
		Timestamp: s.Timestamp.Format("2006-01-02 15:04"),
		URL:       s.URL,
	}
}
