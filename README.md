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

Before running commands, configure the target Kubernetes cluster and namespace. This saves settings to `~/.ingext/config.yaml`. Profiles are keyed by `cluster:namespace`, allowing multiple profiles for the same cluster with different namespaces.

```bash
# Set your default target
ingext config add --cluster <k8s-cluster> --namespace <app-namespace> --context <kubectlContext>  --provider <eks|aks|gke>

# Example — creates profile keyed as "datalake:ingext"
ingext config add --cluster datalake  --namespace ingext --provider eks --context arn:aws:eks:$Region:$AWSAccount:cluster/datalake

```

You can view your current configuration at any time:

```bash
ingext config view
```

Switch between profiles quickly:

```bash
# Switch by cluster:namespace
ingext config use datalake:ingext

# Switch by namespace only (if unique across clusters)
ingext config use ingext
```

List all configured profiles or delete a profile:

```bash
ingext config list
ingext config delete --cluster <cluster-name> --namespace <namespace>
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
| `--cluster` |  | _none_ | Target Kubernetes cluster (required unless using site config). |
| `--namespace` | `-n` | `ingext` | Namespace of the ingext app. |
| `--site-config` |  | `./site_credentials.json` | Path to site credentials file (bypasses Kubernetes). |
| `--site` |  | _none_ | Site hostname from tokenMap (e.g. `demo.cloud.fluencysecurity.com`). |
| `--log-level` | `-l` | `warn` | Log level: `debug`, `info`, `warn`, or `error`. |
| `--version` | `-v` | `false` | Print CLI version (`1.1.0`) and exit. |

### Status (`status`)

Check the current namespace for running services and health checks for core ingext endpoints. Prints a table plus a summary of healthy/degraded/down services.

```bash
ingext status
```

### Authentication (`auth`)

Manage users and access tokens.

```bash
# Users
ingext auth add-user --name foo@gmail.com --role admin --displayName "Foo Bar" --org ingext
ingext auth add-user --name foo@gmail.com --role admin --displayName "Foo Bar" --org ingext --oauth Azure
ingext auth del-user --name foo@gmail.com
ingext auth list-user

# API tokens
ingext auth add-token --name ci-bot --role analyst --description "CI Bot"
ingext auth del-token --name ci-bot
ingext auth list-token

# Site policy
ingext auth set-user-site-policy --username foo@gmail.com --policy my-policy
```

### AI (`ai`)

Register or remove API tokens on a **GPT API** host using explicit URL and bearer token. These commands do **not** use Kubernetes or the main Ingext site from config—you pass `--url` and `--token` on the `ai` command each time.

The same RPC style as `auth add-token` / `del-token` is used (`api/auth` `api_token`). For **register**, the role is fixed to `admin`. CLI flags map to API fields as follows:

| CLI flag | API field |
| --- | --- |
| `--account` | `name` |
| `--display-name` | `description` |

**Flags (on `ingext ai`, shared by subcommands)**

| Flag | Required | Description |
| --- | --- | --- |
| `--url` | Yes | GPT API base URL (e.g. `https://example.com`). |
| `--token` | Yes | Bearer token used to authorize calls to that API. |

**Register**

| Flag | Required | Description |
| --- | --- | --- |
| `--account` | Yes | Account identifier (sent as `name`). |
| `--display-name` | No | Display label (sent as `description`). |

On success, the new API token string is printed to **stdout** (logs go to stderr).

**Unregister**

| Flag | Required | Description |
| --- | --- | --- |
| `--account` | Yes | Account to remove (same `name` as used at registration). |

```bash
# Create a token (role admin); token value is printed on stdout
ingext ai register \
  --url 'https://gpt.example.com' \
  --token "$GPT_BEARER_TOKEN" \
  --account my-service-account \
  --display-name "My integration"

# Remove the token for that account
ingext ai unregister \
  --url 'https://gpt.example.com' \
  --token "$GPT_BEARER_TOKEN" \
  --account my-service-account
```

### Streams (`stream`)

Manage data pipelines (sources, sinks, routers).

```bash
# Sources
ingext stream add-source --name clickstream-v1 --source-type plugin --integration-id <integration-id>
ingext stream add-source --name hec-ingest --source-type hec
ingext stream list-source
ingext stream del-source --id <source-id>

# Sinks
ingext stream add-sink --name datalake-out --sink-type datalake --datalake managed --index <index-name>
ingext stream add-sink --name hec-out --sink-type hec --url https://hec.example --token <token>
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
ingext processor add --name filter-logic --content @./scripts/filter.js [--type fpl_processor] [--desc "Filter logic"]
cat ./scripts/transform.js | ingext processor add --name transform-logic --content -
ingext processor list
ingext processor del --name filter-logic
```

### Integrations (`integration`)

Manage third-party connections.

```bash
ingext integration add --integration slack --name alert-bot --description "Send alerts to Slack" \
  --config key1=value1 --config-bool enabled=true --config-int retries=3 \
  --config-json 'tags=["a","b"]' --secret api_key=xxx --add-source
ingext integration list
ingext integration del --id <integration-id>
```

### KQL Search (`kql`)

Run KQL queries against the datalake and display results as formatted tables.

```bash
# Inline query
ingext kql "MyTable | where status == 200 | take 10"

# Query from file
ingext kql @query.kql

# Save the full JSON response to a file
ingext kql "MyTable | summarize count() by src" --output result.json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--output` | _none_ | Save the full JSON response to a file. |

### Data Lake (`datalake`)

Manage datalakes, indexes, and schemas.

```bash
# Datalakes
ingext datalake add --datalake my-datalake --managed --integration <integration-id>
ingext datalake list

# Indexes
ingext datalake add-index --datalake my-datalake --index events --schema "ingext default"
ingext datalake list-index --datalake my-datalake
ingext datalake del-index --datalake my-datalake --index events

# Schemas
ingext datalake add-schema --name my-schema --schema @./schema.json [--description "My schema"]
ingext datalake update-schema --name my-schema --schema @./schema.json [--description "Updated schema"]
ingext datalake list-schema
ingext datalake delete-schema --name my-schema
```

### Syslog (`syslog`)

Manage syslog ingestion configuration. Supported port types: `tcp`, `udp`, `tls`, `tls-rfc6587`.

```bash
# Get current syslog configuration
ingext syslog get

# Register syslog with specific port types
ingext syslog register --port tcp --port udp

# Update syslog port types
ingext syslog update --port tcp --port tls

# Delete syslog configuration
ingext syslog delete
```

### Grid (`grid`)

Manage grid SaaS accounts.

```bash
# List all grid accounts
ingext grid list-account

# Add a SaaS account
ingext grid add-account --name my-account --region us-east-1 --cluster prod \
  --site-url https://example.fluencysecurity.com --token <token> --display-name "My Account"

# Delete a SaaS account
ingext grid delete-account --name my-account
```

### Collector (`collector`)

List and check status of collectors.

```bash
# List all collectors
ingext collector list

# Get collector status
ingext collector status --collector my-collector
```

### EventWatch (`eventwatch`)

Search event watch summaries, timelines, and rules.

```bash
# Summary search (default: last 1 hour)
ingext eventwatch search_summary --query "keyword"
ingext eventwatch search_summary --query "keyword" --from <unix-ms> --to <unix-ms>

# Timeline search (default: last 1 hour)
ingext eventwatch search_timeline --query "keyword"

# Rule search (no time range)
ingext eventwatch search_rule --query "keyword"
```

### FPL (`fpl`)

Run FPL v2 reports and retrieve task results.

```bash
# Run an FPL report from a JSON file
ingext fpl run --file report.json

# Get task status by ID
ingext fpl get --id 123

# Get task results by ID
ingext fpl results --id 123
```

### Resource (`resource`)

Search resources by type and customer.

```bash
ingext resource --resource-type office365User --customer _all_
```

### EKS Pod Identity Roles (`eks`)

```bash
ingext eks add-assumed-role --name ingest-role --roleArn <role-arn> [--externalId <external-id>]
ingext eks list-assumed-role
ingext eks del-assumed-role --id <role-id>
ingext eks get-pod-role
ingext eks test-assumed-role --roleArn <role-arn> [--externalId <external-id>]
```

### Applications (`application`)

```bash
# List and manage templates
ingext application list
ingext application add --content @./template.yaml
ingext application update --app <template> --content @./template.yaml
ingext application del --app <template>

# Install and manage instances
ingext application install --app <template> --instance <instance> --displayName "My App" --config key=value --secret secretKey=value
ingext application uninstall --app <template> --instance <instance>
ingext application get-instance --app <template> --instance <instance>
```

### Import (`import`)

Import resources from a GitHub repository.

```bash
ingext import processor --type fpl_processor
ingext import application
ingext import schema
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
