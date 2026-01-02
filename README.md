# S3 Bucket Mirror

A high-performance, concurrent Go application for mirroring files between S3 and MinIO buckets with YAML-based configuration.

## Features

- 🚀 **Concurrent Processing**: Configurable worker pool for parallel transfers
- 🔄 **Smart Copying**: Skips files that already exist with matching ETags
- 🌍 **Multi-Storage Support**: Works with AWS S3, MinIO, and hybrid configurations
- 🎯 **Prefix Filtering**: Mirror only specific paths within buckets
- 📝 **YAML Configuration**: Clean, maintainable configuration files
- 🏃 **Dry Run Mode**: Preview operations without actual copying
- 📊 **Progress Tracking**: Verbose logging and summary statistics
- 🔌 **Pluggable Architecture**: Easy to extend with new storage backends

## Installation

```bash
go install github.com/michaelahli/s3-mirror@latest
```

Or build from source:

```bash
git clone https://github.com/michaelahli/s3-mirror.git
cd s3-mirror
go build -o s3-mirror ./cmd/s3-mirror
```

## Prerequisites

- Go 1.21 or later
- For S3: AWS credentials configured (via `~/.aws/credentials`, environment variables, or IAM role)
- For MinIO: Access credentials (specified in config file)

## Configuration

Create a `config.yaml` file with your source and target configuration:

### S3 to S3

```yaml
source:
  type: s3
  bucket: my-source-bucket
  region: us-east-1
  prefix: "data/"  # Optional

target:
  type: s3
  bucket: my-target-bucket
  region: us-west-2
  prefix: "backup/"  # Optional

workers: 20
dry_run: false
verbose: true
```

### MinIO to MinIO

```yaml
source:
  type: minio
  bucket: source-bucket
  endpoint: minio1.example.com:9000
  access_key_id: minioadmin
  secret_access_key: minioadmin
  use_ssl: false

target:
  type: minio
  bucket: target-bucket
  endpoint: minio2.example.com:9000
  access_key_id: minioadmin
  secret_access_key: minioadmin
  use_ssl: false

workers: 10
dry_run: false
verbose: false
```

### S3 to MinIO (Hybrid)

```yaml
source:
  type: s3
  bucket: aws-bucket
  region: us-east-1

target:
  type: minio
  bucket: minio-bucket
  endpoint: minio.example.com:9000
  access_key_id: your-access-key
  secret_access_key: your-secret-key
  use_ssl: true

workers: 15
```

### MinIO to S3 (Hybrid)

```yaml
source:
  type: minio
  bucket: minio-bucket
  endpoint: localhost:9000
  access_key_id: minioadmin
  secret_access_key: minioadmin
  use_ssl: false

target:
  type: s3
  bucket: aws-archive
  region: us-east-1

workers: 10
```

## Usage

### Basic Usage

```bash
s3-mirror -config config.yaml
```

### Command Line Options

```bash
s3-mirror \
  -config config.yaml \
  -workers 20 \
  -verbose \
  -dry-run
```

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-config` | Path to YAML configuration file | `config.yaml` |
| `-workers` | Number of concurrent workers (overrides config) | From config (default: 10) |
| `-dry-run` | Preview operations without copying | `false` |
| `-verbose` | Enable detailed logging | `false` |

## Configuration Reference

### Source/Target Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | Storage type: `s3` or `minio` |
| `bucket` | string | Yes | Bucket name |
| `region` | string | S3 only | AWS region (default: `us-east-1`) |
| `endpoint` | string | MinIO only | MinIO endpoint (e.g., `localhost:9000`) |
| `access_key_id` | string | MinIO only | MinIO access key |
| `secret_access_key` | string | MinIO only | MinIO secret key |
| `use_ssl` | boolean | MinIO only | Use HTTPS for MinIO (default: `false`) |
| `prefix` | string | No | Filter/prefix for objects |

### Global Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `workers` | integer | Number of concurrent workers | `10` |
| `dry_run` | boolean | Preview mode without copying | `false` |
| `verbose` | boolean | Enable detailed logging | `false` |

## How It Works

1. **Load Configuration**: Reads YAML config file with source/target settings
2. **Initialize Clients**: Creates appropriate storage clients (S3/MinIO) based on config
3. **List Objects**: Retrieves all objects from source bucket (with optional prefix filter)
4. **Check Existing**: For each object, checks if it exists in target with matching ETag
5. **Copy Objects**: Uses concurrent workers to copy objects via download/upload
6. **Report Results**: Provides summary statistics

## Use Cases

### Backup AWS S3 to MinIO

Perfect for creating local backups of cloud data:

```yaml
source:
  type: s3
  bucket: production-data
  region: us-east-1

target:
  type: minio
  bucket: backup
  endpoint: backup.local:9000
  access_key_id: admin
  secret_access_key: password
  use_ssl: false
```

### Migrate Between MinIO Clusters

Ideal for data migration between MinIO deployments:

```yaml
source:
  type: minio
  bucket: old-cluster-data
  endpoint: old.minio.local:9000
  # ... credentials ...

target:
  type: minio
  bucket: new-cluster-data
  endpoint: new.minio.local:9000
  # ... credentials ...
```

### Archive to S3 from MinIO

Move cold data from MinIO to S3 for long-term storage:

```yaml
source:
  type: minio
  bucket: active-data
  endpoint: minio.local:9000
  # ... credentials ...
  prefix: "2023/"  # Only archive 2023 data

target:
  type: s3
  bucket: archive-bucket
  region: us-west-2
```

## MinIO Setup

To test with MinIO locally:

```bash
# Run MinIO with Docker
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"

# Access MinIO Console at http://localhost:9001
```

## Performance Tips

- **Worker Count**: Start with 10-20 workers and adjust based on performance
- **Network**: Run from server close to storage for better throughput
- **Prefix Filtering**: Use prefixes to mirror specific subsets of data
- **Dry Run**: Always test with `-dry-run` first to preview operations
- **Verbose Mode**: Use `-verbose` to monitor progress in real-time

## Examples

### Mirror with prefix filter

```bash
s3-mirror -config config.yaml -verbose
```

### Dry run to preview

```bash
s3-mirror -config config.yaml -dry-run
```

### Override worker count

```bash
s3-mirror -config config.yaml -workers 50
```

## Security Best Practices

1. **AWS Credentials**: Use IAM roles when running on EC2
2. **MinIO Credentials**: Store credentials securely, use environment variables if needed
3. **SSL/TLS**: Enable `use_ssl: true` for MinIO in production
4. **Access Control**: Grant minimum required permissions:
   - Source: `s3:ListBucket`, `s3:GetObject`
   - Target: `s3:ListBucket`, `s3:PutObject`, `s3:GetObject`

## Troubleshooting

### Connection Issues

- For MinIO: Verify endpoint is accessible and credentials are correct
- For S3: Check AWS credentials and region settings

### Performance Issues

- Increase worker count for better throughput
- Check network bandwidth and latency
- Monitor CPU and memory usage

### Permission Errors

- Verify IAM policies for S3
- Check MinIO user policies and access rights

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [x] S3 support
- [x] MinIO support
- [x] YAML configuration
- [x] Hybrid S3/MinIO support
- [ ] Resume capability for interrupted transfers
- [ ] Incremental sync based on timestamp
- [ ] Progress bar for visual feedback
- [ ] Delete mode (sync with deletions)
- [ ] Support for custom storage classes
- [ ] Bandwidth throttling
- [ ] Multi-part upload for large files
- [ ] Encryption support
- [ ] Compression options

## Acknowledgments

Built with:

- [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2)
- [MinIO Go SDK](https://github.com/minio/minio-go)
