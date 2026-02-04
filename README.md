# gcsgate

A Google App Engine proxy for serving Google Cloud Storage files behind IAP (Identity-Aware Proxy) authentication.

## Overview

GCS buckets cannot be directly protected by IAP. This proxy runs on App Engine with IAP enabled, providing authenticated access to your private GCS files.

```
User -> IAP -> App Engine (gcsgate) -> GCS
```

## Deployment

```bash
make deploy PROJECT=your-gcp-project SERVICE_ACCOUNT=gcsgate@your-gcp-project.iam.gserviceaccount.com
```

## Setup

### 1. Create a Service Account

Create a dedicated service account for gcsgate:

```bash
gcloud iam service-accounts create gcsgate \
  --project=your-gcp-project \
  --display-name="gcsgate"
```

### 2. Grant GCS Access

Grant the service account read access to specific buckets:

```bash
gcloud storage buckets add-iam-policy-binding gs://your-bucket \
  --member="serviceAccount:gcsgate@your-gcp-project.iam.gserviceaccount.com" \
  --role="roles/storage.objectViewer"
```

Repeat for each bucket you want to serve.

### 3. Enable IAP

1. Go to [IAP settings](https://console.cloud.google.com/security/iap) in Google Cloud Console
2. Enable IAP for your App Engine application
3. Grant `IAP-secured Web App User` role to users/groups who need access

## Usage

URL paths map directly to GCS objects:

| URL Path | GCS Object |
|----------|------------|
| `/{bucket}/{path}` | `gs://{bucket}/{path}` |

Examples:

| Request | GCS Object |
|---------|------------|
| `/my-bucket/reports/2024/data.html` | `gs://my-bucket/reports/2024/data.html` |
| `/my-bucket/images/chart.png` | `gs://my-bucket/images/chart.png` |
| `/another-bucket/dir/file.csv` | `gs://another-bucket/dir/file.csv` |

Full URL: `https://your-gcp-project.appspot.com/{bucket}/{path}`

Relative path references in HTML files (`./image.png`, `../style.css`, etc.) work correctly.

## Configuration

See `app.yaml` for configuration options. Key settings:

- `service`: Service name (default: `gcsgate`)
- `service_account`: Service account for GCS access
- `instance_class`: Instance size (default: `F1`)
- `automatic_scaling`: Scaling configuration

## Local Development

```bash
# Build
make build

# Run
make run
```

Note: Local requests return 401 Unauthorized because the IAP header is missing. For testing, add the `x-goog-iap-jwt-assertion` header to your requests.

## License

MIT
