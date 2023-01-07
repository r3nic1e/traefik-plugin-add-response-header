# traefik-plugin-add-response-header

Traefik proxy plugin to copy request header to response.
I use this plugin to copy assigned trace ID to response headers.

**This plugin copies full response body to temporary buffer. Serious performance issues may occur**

## Configuration

Add plugin:
```yaml
experimental:
  plugins:
    add-response-header:
      moduleName: github.com/r3nic1e/traefik-plugin-add-response-header
      version: v0.4.0
```

Configure middleware:
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: add-response-header
spec:
  plugin:
    add-response-header:
      from: "X-B3-TraceId"
      to: "X-B3-TraceId"
```

You can also use regexp:
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: add-response-header
spec:
  plugin:
    add-response-header:
      from: "X-B3-TraceId"
      to: "X-B3-TraceId"
      overwrite: true
      regexp: "^(.*)/(.*)$"
      replacement: "$1:$2"
```