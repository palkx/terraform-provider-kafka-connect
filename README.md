# `terraform-plugin-kafka-connect`

[![test](https://github.com/Mongey/terraform-provider-kafka-connect/actions/workflows/test.yml/badge.svg)](https://github.com/Mongey/terraform-provider-kafka-connect/actions/workflows/test.yml)

A [Terraform][1] plugin for managing [Apache Kafka Connect][2].

## Installation

Download and extract the [latest
release](https://github.com/Mongey/terraform-provider-kafka-connect/releases/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`)

## Example

Configure the provider directly, or set the ENV variable `KAFKA_CONNECT_URL`

```hcl
provider "kafka-connect" {
  url = "http://localhost:8083"
  basic_auth_username = "user" # Optional
  basic_auth_password = "password" # Optional

  # For TLS
  tls_auth_crt = "/tmp/cert.pem" # Optional
  tls_auth_key = "/tmp/key.pem " # Optional
  tls_auth_is_insecure = true   # Optionnal if you do not want to check CA
}

resource "kafka-connect_connector" "sqlite-sink" {
  name = "sqlite-sink"

  config = {
    "name"            = "sqlite-sink"
    "connector.class" = "io.confluent.connect.jdbc.JdbcSinkConnector"
    "tasks.max"       = 1
    "topics"          = "orders"
    "connection.url"  = "jdbc:sqlite:test.db"
    "auto.create"     = "true"
    "connection.user" = "admin"
  }

  config_sensitive = {
    "connection.password" = "this-should-never-appear-unmasked"
  }
}
```

## Provider Properties

| Property              | Type              | Example                 | Alternative environment variable name |
| --------------------- | ----------------- | ----------------------- | ------------------------------------- |
| `url`                 | URL               | "http://localhost:8083" | `KAFKA_CONNECT_URL`                   |
| `headers`             | Map[String]String | {foo = "bar"}           | N/A                                   |
| `basic_auth_username` | String            | "user"                  | `KAFKA_CONNECT_BASIC_AUTH_USERNAME`   |
| `basic_auth_password` | String            | "password"              | `KAFKA_CONNECT_BASIC_AUTH_PASSWORD`   |
| `tls_skip_insecure`   | Bool              | false                   | `KAFKA_CONNECT_TLS_IS_INSECURE`       |
| `tls_root_ca_path`    | String            | "/path/to/ca-cert.pem"  | `KAFKA_CONNECT_SSL_ROOT_CA_FILE_PATH` |
| `tls_auth_crt_path`   | String            | "/path/to/cert.pem"     | `KAFKA_CONNECT_TLS_AUTH_CRT_PATH`     |
| `tls_auth_key_path`   | String            | "/path/to/key.pem"      | `KAFKA_CONNECT_TLS_AUTH_KEY_PATH`     |

## Resource Properties

| Property           | Type      | Description                                                  |
| ------------------ | --------- | ------------------------------------------------------------ |
| `name`             | String    | Connector name                                               |
| `config`           | HCL Block | Connector configuration                                      |
| `config_sensitive` | HCL Block | Sensitive connector configuration. Will be masked in output. |

## Developing

0. [Install go][install-go]
1. Clone repository to: `$GOPATH/src/github.com/Mongey/terraform-provider-kafka-connect`
   ```bash
   mkdir -p $GOPATH/src/github.com/Mongey/terraform-provider-kafka-connect; cd $GOPATH/src/github.com/Mongey/
   git clone https://github.com/Mongey/terraform-provider-kafka-connect.git
   ```
2. Build the provider `make build`
3. Run the tests `make test`
4. Run acceptance tests `make testacc`

One time compilation with docker:

1. `env UID=$(id -u) GID=$(id -g) GOOS=$(uname -o | tr '[:upper:]' '[:lower:]') GOARCH=$(uname -m) docker compose -f local-build.yml up --build`

[1]: https://www.terraform.io
[2]: https://kafka.apache.org/documentation/#connect
[third-party-plugins]: https://www.terraform.io/docs/configuration/providers.html#third-party-plugins
[install-go]: https://golang.org/doc/install#install
