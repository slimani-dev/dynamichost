# DynamicHost Plugin for Traefik

[![Build Status](https://github.com/yourusername/dynamichost-plugin/workflows/Main/badge.svg?branch=master)](https://github.com/yourusername/dynamichost-plugin/actions)

The `DynamicHost` plugin is a middleware for [Traefik](https://traefik.io) that dynamically rewrites the host header based on a regex pattern. This allows flexible host transformation based on request attributes.

## Features
- Uses a regex pattern to match and transform the host dynamically.
- Allows custom host structures for different incoming requests.
- Fully compatible with Traefik's middleware system.

## Configuration

To use the `DynamicHost` plugin, you must define it in the **static configuration** of Traefik:

```yaml
# Static configuration
experimental:
  plugins:
    dynamichost:
      moduleName: github.com/yourusername/dynamichost-plugin
      version: v0.1.0
```

Then, you can configure it dynamically:

```yaml
# Dynamic configuration
http:
  routers:
    my-router:
      rule: Host(`example.localhost`)
      service: my-service
      entryPoints:
        - web
      middlewares:
        - dynamic-host

  services:
    my-service:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:8080

  middlewares:
    dynamic-host:
      plugin:
        dynamichost:
          regexPattern: "^([^.]+)\\.localhost$"
          newHost: "$1.example.com"
```

### Parameters

| Parameter       | Type   | Description |
|----------------|--------|-------------|
| `regexPattern` | string | The regex pattern used to match the original host. |
| `newHost`      | string | The new host format using regex capture groups. |

### Example Behavior

| Incoming Host      | Transformed Host  |
|-------------------|------------------|
| `abc.localhost`   | `abc.example.com` |
| `test.localhost`  | `test.example.com` |

## Development

To develop locally, use the local plugins feature of Traefik:

```yaml
# Static configuration
experimental:
  localPlugins:
    dynamichost:
      moduleName: github.com/yourusername/dynamichost-plugin
```

Place your plugin source code in the `./plugins-local` directory:

```
./plugins-local/
    └── src
        └── github.com
            └── yourusername
                └── dynamichost-plugin
                    ├── main.go
                    ├── go.mod
                    ├── LICENSE
                    ├── README.md
```

## Logs

Currently, logs can be written using `os.Stdout.WriteString("...")` or `os.Stderr.WriteString("...")`. Future versions may support structured logging.

## Plugin Catalog

To add this plugin to the Traefik Plugin Catalog, ensure:

- The repository has the `traefik-plugin` topic.
- A `.traefik.yml` manifest is present with valid configuration.

Example `.traefik.yml`:

```yaml
displayName: DynamicHost
summary: Middleware that rewrites host dynamically based on regex.
type: middleware
import: github.com/yourusername/dynamichost-plugin
testData:
  regexPattern: "^([^.]+)\\.localhost$"
  newHost: "$1.example.com"
```

## License
This project is licensed under the MIT License.

