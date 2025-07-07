package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/taraxa/snapshots-api/internal/models"
)

// SnapshotParser handles parsing of snapshot filenames
type SnapshotParser struct {
	// Regex pattern for snapshot filename: <network>-<full/light>-db-block-<blocknumber>-<timestamp>.tar.gz
	pattern *regexp.Regexp
}

// NewSnapshotParser creates a new snapshot parser
func NewSnapshotParser() *SnapshotParser {
	// Pattern matches: network-type-db-block-blocknumber-timestamp.tar.gz
	pattern := regexp.MustCompile(`^(mainnet|testnet|devnet)-(full|light)-db-block-(\d+)-(\d{8}-\d{6})\.tar\.gz$`)
	return &SnapshotParser{
		pattern: pattern,
	}
}

// ParseSnapshot parses a snapshot filename and returns a Snapshot struct
func (p *SnapshotParser) ParseSnapshot(filename, baseURL string) (*models.Snapshot, error) {
	matches := p.pattern.FindStringSubmatch(filename)
	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid snapshot filename format: %s", filename)
	}

	network := models.Network(matches[1])
	snapshotType := models.SnapshotType(matches[2])
	blockStr := matches[3]
	timestampStr := matches[4]

	// Parse block number
	block, err := strconv.ParseInt(blockStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid block number %s: %w", blockStr, err)
	}

	// Parse timestamp (format: YYYYMMDD-HHMMSS)
	timestamp, err := time.Parse("20060102-150405", timestampStr)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %s: %w", timestampStr, err)
	}

	// Construct public URL
	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(baseURL, "/o"), filename)

	return &models.Snapshot{
		Network:   network,
		Type:      snapshotType,
		Block:     block,
		Timestamp: timestamp,
		URL:       url,
		Filename:  filename,
	}, nil
}

// IsValidNetwork checks if the network is supported
func (p *SnapshotParser) IsValidNetwork(network string) bool {
	switch models.Network(network) {
	case models.NetworkMainnet, models.NetworkTestnet, models.NetworkDevnet:
		return true
	default:
		return false
	}
}
