# traefik-plugin-add-response-header

Traefik proxy plugin to copy request header to response.
I use this plugin to copy assigned trace ID to response headers.

## Configuration

Add plugin:
```yaml
experimental:
  plugins:
    add-response-header:
      moduleName: github.com/argyle-engineering/copy-header-value-traefik-plugin
      version: v0.1.0
```

Configure middleware:
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: add-response-header
spec:
  plugin:
    copy-header-value:
      from: "X-B3-TraceId"
      to: "X-B3-TraceId"
      overwrite: false
```