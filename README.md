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
ingext config set --cluster <k8s-cluster> --namespace <app-namespace> --context <kubectlContext>  --provider <eks|aks|gke>

# Example
ingext config set --cluster datalake  --namespace ingext --provider eks --context arn:aws:eks:$Region:$AWSAccount:cluster/datalake 

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

### Global flags

| Flag | Shorthand | Default | Description |
| --- | --- | --- | --- |
| `--cluster` |  | _none_ | Target Kubernetes cluster (required unless set via config). |
| `--namespace` | `-n` | `ingext` | Namespace of the ingext app. |
| `--log-level` | `-l` | `warn` | Log level: `debug`, `info`, `warn`, or `error`. |
| `--version` | `-v` | `false` | Print CLI version (`1.0.0`) and exit. |

### Authentication (`auth`)

Manage users and access tokens.

```bash
# Users
ingext auth add-user --name foo@gmail.com --role admin --displayName "Foo Bar" --org ingext
ingext auth del-user --name foo@gmail.com
ingext auth list-user

# API tokens
ingext auth add-token --name ci-bot --role analyst --displayName "CI Bot"
ingext auth del-token --name ci-bot
ingext auth list-token
```

### Streams (`stream`)

Manage data pipelines (sources, sinks, routers).

```bash
# Sources
ingext stream add-source --name clickstream-v1 --source-type plugin --integration-id <integration-id>
ingext stream add-source --name hec-ingest --source-type hec --url https://hec.example --token <token>
ingext stream list-source
ingext stream del-source --id <source-id>

# Sinks
ingext stream add-sink --name datalake-out --sink-type datalake --integration-id <integration-id>
ingext stream add-sink --name webhook-out --sink-type webhook --url https://example.com/hook
ingext stream list-sink
ingext stream del-sink --id <sink-id>

# Routers and wiring
ingext stream add-router --processor my-processor --router-name main-router
ingext stream connect-router --source-id <source-id> --router-id <router-id>
ingext stream connect-sink --router-id <router-id> --sink-id <sink-id>
```

### Processors (`processor`)

Deploy data processors. Supports piping input via `-` and file loading via `@path`.

```bash
ingext processor add --name filter-logic --content @./scripts/filter.js
cat ./scripts/transform.js | ingext processor add --name transform-logic --content -
ingext processor list
ingext processor del --name filter-logic
```

### Integrations (`integration`)

Manage third-party connections.

```bash
ingext integration add --integration slack --name alert-bot --description "Send alerts to Slack"
ingext integration list
ingext integration del --id <integration-id>
```

### Data Lake (`datalake`)

Manage datalakes and their indexes.

```bash
ingext datalake add --datalake my-datalake --managed --integration <integration-id>
ingext datalake list
ingext datalake add-index --datalake my-datalake --index events --schema "ingext default"
ingext datalake list-index --datalake my-datalake
ingext datalake del-index --datalake my-datalake --index events
```

### EKS Pod Identity Roles (`eks`)

```bash
ingext eks add-assumed-role --name ingest-role --roleArn <role-arn> [--externalId <external-id>]
ingext eks list-assumed-role
ingext eks del-assumed-role --id <role-id>
```

### Applications (`application`)

```bash
ingext application list-template
ingext application install --app <template> --instance <instance> --displayName "My App" --set key=value
ingext application uninstall --app <template> --instance <instance>
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
go get "github.com/google/go-github/v64/github"
go mod tidy

```

## License

[MIT](https://www.google.com/search?q=LICENSE)
