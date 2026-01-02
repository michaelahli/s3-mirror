package mirror

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/michaelahli/s3-mirror/internal/config"
	"github.com/michaelahli/s3-mirror/internal/storage"
)

// Stats holds statistics about the mirror operation
type Stats struct {
	TotalObjects     int64
	CopiedObjects    int64
	SkippedObjects   int64
	Errors           int64
	BytesTransferred int64
}

// Service handles the mirroring logic
type Service struct {
	sourceClient storage.Client
	targetClient storage.Client
	config       *config.Config
	stats        Stats
	ctx          context.Context
}

// NewService creates a new mirror service
func NewService(sourceClient, targetClient storage.Client, cfg *config.Config) *Service {
	return &Service{
		sourceClient: sourceClient,
		targetClient: targetClient,
		config:       cfg,
		ctx:          context.Background(),
	}
}

// Mirror performs the mirroring operation
func (s *Service) Mirror() (*Stats, error) {
	log.Printf("Starting mirror from %s to %s", s.sourceClient.GetBucket(), s.targetClient.GetBucket())
	if s.config.Source.Prefix != "" {
		log.Printf("Using source prefix filter: %s", s.config.Source.Prefix)
	}
	if s.config.DryRun {
		log.Println("DRY RUN MODE - No actual copying will occur")
	}

	// List all objects in source bucket
	objects, err := s.sourceClient.ListObjects(s.ctx, s.config.Source.Prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list source objects: %w", err)
	}

	atomic.StoreInt64(&s.stats.TotalObjects, int64(len(objects)))
	log.Printf("Found %d objects to process", len(objects))

	// Create work channel and worker pool
	objectChan := make(chan storage.ObjectInfo, len(objects))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go s.worker(&wg, objectChan)
	}

	// Send objects to workers
	for _, obj := range objects {
		objectChan <- obj
	}
	close(objectChan)

	// Wait for all workers to complete
	wg.Wait()

	return &s.stats, nil
}

// worker processes objects from the work channel
func (s *Service) worker(wg *sync.WaitGroup, objectChan <-chan storage.ObjectInfo) {
	defer wg.Done()

	for obj := range objectChan {
		if err := s.copyObject(obj); err != nil {
			log.Printf("Error copying %s: %v", obj.Key, err)
			atomic.AddInt64(&s.stats.Errors, 1)
		}
	}
}

// copyObject copies a single object from source to target
func (s *Service) copyObject(obj storage.ObjectInfo) error {
	// Apply target prefix if specified
	targetKey := obj.Key
	if s.config.Target.Prefix != "" {
		targetKey = s.config.Target.Prefix + obj.Key
	}

	// Check if object already exists in target
	targetObj, err := s.targetClient.HeadObject(s.ctx, targetKey)
	if err == nil {
		// Object exists, check if it's the same
		if targetObj.ETag == obj.ETag {
			if s.config.Verbose {
				log.Printf("Skipping %s (already exists with same ETag)", obj.Key)
			}
			atomic.AddInt64(&s.stats.SkippedObjects, 1)
			return nil
		}
	}

	// Perform the copy
	if s.config.DryRun {
		log.Printf("[DRY RUN] Would copy: %s -> %s (%d bytes)", obj.Key, targetKey, obj.Size)
		atomic.AddInt64(&s.stats.CopiedObjects, 1)
		return nil
	}

	// Download from source
	body, err := s.sourceClient.GetObject(s.ctx, obj.Key)
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer body.Close()

	// Upload to target
	err = s.targetClient.PutObject(s.ctx, targetKey, body, obj.Size)
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	if s.config.Verbose {
		log.Printf("Copied: %s -> %s (%d bytes)", obj.Key, targetKey, obj.Size)
	}

	atomic.AddInt64(&s.stats.CopiedObjects, 1)
	atomic.AddInt64(&s.stats.BytesTransferred, obj.Size)

	return nil
}
