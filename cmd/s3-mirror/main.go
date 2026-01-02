package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/michaelahli/s3-mirror/internal/config"
	"github.com/michaelahli/s3-mirror/internal/mirror"
	"github.com/michaelahli/s3-mirror/internal/storage"
)

func main() {
	// Command line flags
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	workers := flag.Int("workers", 0, "Number of concurrent workers (overrides config)")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run without copying")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	// Load configuration from YAML
	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override with command line flags if provided
	if *workers > 0 {
		cfg.Workers = *workers
	}
	if *dryRun {
		cfg.DryRun = true
	}
	if *verbose {
		cfg.Verbose = true
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize storage clients
	sourceClient, err := storage.NewClient(cfg.Source)
	if err != nil {
		log.Fatalf("Failed to create source storage client: %v", err)
	}

	targetClient, err := storage.NewClient(cfg.Target)
	if err != nil {
		log.Fatalf("Failed to create target storage client: %v", err)
	}

	// Create mirror service
	mirrorService := mirror.NewService(sourceClient, targetClient, cfg)

	// Run the mirror operation
	stats, err := mirrorService.Mirror()
	if err != nil {
		log.Fatalf("Mirror operation failed: %v", err)
	}

	// Print summary
	fmt.Println("\n=== Mirror Summary ===")
	fmt.Printf("Total objects processed: %d\n", stats.TotalObjects)
	fmt.Printf("Objects copied: %d\n", stats.CopiedObjects)
	fmt.Printf("Objects skipped: %d\n", stats.SkippedObjects)
	fmt.Printf("Errors: %d\n", stats.Errors)
	fmt.Printf("Total bytes transferred: %d (%.2f MB)\n", stats.BytesTransferred, float64(stats.BytesTransferred)/(1024*1024))

	if stats.Errors > 0 {
		os.Exit(1)
	}
}
