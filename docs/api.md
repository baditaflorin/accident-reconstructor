# API

OpenAPI contract:

https://github.com/baditaflorin/accident-reconstructor/blob/main/api/openapi.yaml

Local base URL:

http://localhost:8080

Health:

```sh
curl http://localhost:8080/healthz
```

Toolchain:

```sh
curl http://localhost:8080/api/v1/tools
```

Create a case:

```sh
curl -F case_name=Demo -F scale_meters=10 -F videos=@dashcam.mp4 http://localhost:8080/api/v1/cases
```

Read artifact:

```sh
curl http://localhost:8080/api/v1/cases/CASE_ID/artifact
```
