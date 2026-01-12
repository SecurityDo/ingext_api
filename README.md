# ingext CLI

`ingext` is a command-line interface tool for managing Ingext resources on Kubernetes. It allows users to manage streams, processors, integrations, and authentication through a structured, AWS-CLI-style interface.

## Features

* **Standardized CLI:** Follows the intuitive `noun verb [flags]` pattern (e.g., `ingext stream add source`).
* **Kubernetes Native:** Connects directly to your clusters using your local `kubeconfig` context.
* **Pipe Friendly:** Designed for automation—strictly separates data output (STDOUT) from logs (STDERR) and supports reading files from STDIN.
* **Smart Config:** Hierarchical configuration (Flags > Env Vars > Config File).

## Installation

### From Source

Requirements: Go 1.21+

```bash
# Clone the repository
git clone https://github.com/your-org/ingext.git
cd ingext

# Build the binary
go build -o ingext cmd/ingext/main.go

# (Optional) Move to path
install -m 755 ingext /usr/local/bin/.

```

## Configuration

Before running commands, configure the target Kubernetes cluster and namespace. This saves settings to `~/.ingext/config.yaml`.

```bash
# Set your default target
ingext config --cluster <k8s-cluster> --namespace <app-namespace> --context <kubectlContext>  --provider <eks|aks|gke>

# Example
ingext config --cluster datalake  --namespace ingext --provider eks --context arn:aws:eks:$Region:$AWSAccount:cluster/datalake 

```

You can view your current configuration at any time:

```bash
ingext config view

```

**Environment Variables**
You can override defaults using `INGEXT_` prefixed variables:

```bash
export INGEXT_CLUSTER=prod-cluster
export INGEXT_NAMESPACE=ingext

```

## Usage

### 1. Authentication (`auth`)

Manage users and access tokens.

```bash
# Add a new user
ingext user add --name foo@gmail.com --role admin --displayName "Foo Bar"

# Add an API token for a user
ingext user del --name foo@gmail.com

```

### 2. Streams (`stream`)

Manage data pipelines (Sources, Sinks, Routers).

```bash
# Add a stream source
ingext stream add source --name clickstream-v1

# Add a stream sink
ingext stream add sink --name s3-archive

```

### 3. Processors (`processor`)

Deploy data processors. This command supports piping input via `-`.

```bash
# Deploy from a local file
ingext processor add --name filter-logic --file ./scripts/filter.js

# Deploy from a pipe (stdin)
cat ./scripts/transform.js | ingext processor add --name transform-logic --file -

```

### 4. Integrations (`integration`)

Manage third-party connections.

```bash
ingext integration add --integration slack --name alert-bot
ingext integration del --name alert-bot

```

### 5. Data Lake (`lake`)

Manage storage indexing.

```bash
ingext lake add index \
  --storage s3 \
  --bucket my-datalake \
  --prefix /events/raw \
  --storageaccount my-account

```

## Development

### Project Structure

The project follows the Standard Go Project Layout:

| Path | Description |
| --- | --- |
| `cmd/ingext/` | Application entry point (`main.go`). |
| `internal/commands/` | Cobra command definitions and flag parsing. |
| `internal/api/` | Business logic and Kubernetes client (`client-go`). |
| `internal/config/` | Configuration loading (Viper). |

### Kubernetes Dependency Note

This project uses `client-go` v0.35.0. If you change versions, ensure all k8s libraries match exactly to avoid build errors:

```bash
go get k8s.io/client-go@v0.35.0 k8s.io/api@v0.35.0 k8s.io/apimachinery@v0.35.0
go mod tidy

```

## License

[MIT](https://www.google.com/search?q=LICENSE)
