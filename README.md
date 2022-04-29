# Healthy but not Ready

## Debug liveness and readiness probes in Kubernetes

This app enables you to configure different properties for readiness & liveness probes, and gracefull shutdown period. Useful to debug issues on scaling out or in on Kubernetes.

Environment variables:

* `DELAY_UNTIL_READY` - time until readiness probes pass. Default value `15s`
* `GRACEFUL_SHUTDOWN_DURATION` - duration for graceful shutdown. Default value `1ms` 
* `DELAY_REQUESTS_DURATION` - delay duration for responses to clients

## Endpoints

The app has three endpoints:
- `/healthy` - used for liveness probes
- `/ready` - used for readiness probes
- `/hello` - regular request that delays the response according to the environment variable `DELAY_REQUESTS_DURATION`

## Running the app

```
go mod download 
go run main.go
```

## Running in kubernetes

Check the file [deployment.yaml](./deployment.yaml)
