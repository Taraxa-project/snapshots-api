package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/taraxa/snapshots-api/internal/models"
	"github.com/taraxa/snapshots-api/internal/parser"
)

// GCPStorageResponse represents the response from GCP Storage API
type GCPStorageResponse struct {
	Kind  string `json:"kind"`
	Items []struct {
		Name string `json:"name"`
	} `json:"items"`
}

// SnapshotService handles snapshot operations
type SnapshotService struct {
	bucketName string
	bucketURL  string
	parser     *parser.SnapshotParser
	cache      map[models.Network]*models.NetworkSnapshots
	cacheTime  time.Time
	mutex      sync.RWMutex
	cacheTTL   time.Duration
}

// NewSnapshotService creates a new snapshot service
func NewSnapshotService(bucketName, bucketURL string) *SnapshotService {
	return &SnapshotService{
		bucketName: bucketName,
		bucketURL:  bucketURL,
		parser:     parser.NewSnapshotParser(),
		cache:      make(map[models.Network]*models.NetworkSnapshots),
		cacheTTL:   5 * time.Minute, // Cache for 5 minutes
	}
}

// GetSnapshots retrieves snapshots for a specific network (backward compatibility)
func (s *SnapshotService) GetSnapshots(network models.Network) (*models.NetworkSnapshots, error) {
	return s.GetSnapshotsWithAuth(network, true)
}

// GetSnapshotsWithAuth retrieves snapshots for a specific network with authentication filtering
func (s *SnapshotService) GetSnapshotsWithAuth(network models.Network, authenticated bool) (*models.NetworkSnapshots, error) {
	s.mutex.RLock()
	cached, exists := s.cache[network]
	cacheValid := time.Since(s.cacheTime) < s.cacheTTL
	s.mutex.RUnlock()

	if exists && cacheValid {
		// If not authenticated, filter out full snapshots from cached data
		if !authenticated {
			filteredResult := &models.NetworkSnapshots{
				Light: cached.Light,
				// Full is omitted (nil) for unauthenticated requests
			}
			return filteredResult, nil
		}
		return cached, nil
	}

	// Fetch fresh data
	snapshots, err := s.fetchSnapshots()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch snapshots: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Update cache
	s.cache = s.processSnapshots(snapshots)
	s.cacheTime = time.Now()

	result, exists := s.cache[network]
	if !exists {
		return &models.NetworkSnapshots{}, nil
	}

	// If not authenticated, filter out full snapshots
	if !authenticated {
		filteredResult := &models.NetworkSnapshots{
			Light: result.Light,
			// Full is omitted (nil) for unauthenticated requests
		}
		return filteredResult, nil
	}

	return result, nil
}

// fetchSnapshots retrieves all snapshots from GCP bucket
func (s *SnapshotService) fetchSnapshots() ([]*models.Snapshot, error) {
	resp, err := http.Get(s.bucketURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bucket contents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GCP API returned status %d", resp.StatusCode)
	}

	var gcpResp GCPStorageResponse
	if err := json.NewDecoder(resp.Body).Decode(&gcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode GCP response: %w", err)
	}

	var snapshots []*models.Snapshot
	baseURL := fmt.Sprintf("https://storage.googleapis.com/%s", s.bucketName)

	for _, item := range gcpResp.Items {
		snapshot, err := s.parser.ParseSnapshot(item.Name, baseURL)
		if err != nil {
			// Skip invalid filenames (not all files in bucket are snapshots)
			continue
		}
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// processSnapshots groups snapshots by network and finds the latest for each type
func (s *SnapshotService) processSnapshots(snapshots []*models.Snapshot) map[models.Network]*models.NetworkSnapshots {
	result := make(map[models.Network]*models.NetworkSnapshots)

	// Group by network and type
	networkSnapshots := make(map[models.Network]map[models.SnapshotType][]*models.Snapshot)

	for _, snapshot := range snapshots {
		if _, exists := networkSnapshots[snapshot.Network]; !exists {
			networkSnapshots[snapshot.Network] = make(map[models.SnapshotType][]*models.Snapshot)
		}
		networkSnapshots[snapshot.Network][snapshot.Type] = append(
			networkSnapshots[snapshot.Network][snapshot.Type],
			snapshot,
		)
	}

	// Find latest snapshot for each network and type
	for network, typeSnapshots := range networkSnapshots {
		networkResult := &models.NetworkSnapshots{}

		for snapshotType, snapshots := range typeSnapshots {
			latest := s.findLatestSnapshot(snapshots)
			if latest != nil {
				switch snapshotType {
				case models.SnapshotTypeFull:
					networkResult.Full = latest.ToSnapshotInfo()
				case models.SnapshotTypeLight:
					networkResult.Light = latest.ToSnapshotInfo()
				}
			}
		}

		result[network] = networkResult
	}

	return result
}

// findLatestSnapshot finds the snapshot with the highest block number, or latest timestamp if blocks are equal
func (s *SnapshotService) findLatestSnapshot(snapshots []*models.Snapshot) *models.Snapshot {
	if len(snapshots) == 0 {
		return nil
	}

	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].Block != snapshots[j].Block {
			return snapshots[i].Block > snapshots[j].Block
		}
		return snapshots[i].Timestamp.After(snapshots[j].Timestamp)
	})

	return snapshots[0]
}

// IsValidNetwork checks if a network string is valid
func (s *SnapshotService) IsValidNetwork(network string) bool {
	return s.parser.IsValidNetwork(network)
}

// GetAllNetworks returns all available networks
func (s *SnapshotService) GetAllNetworks() []models.Network {
	return []models.Network{
		models.NetworkMainnet,
		models.NetworkTestnet,
		models.NetworkDevnet,
	}
}
