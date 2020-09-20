# How to develop Go gRPC microservices and deploy in Kubernetes

## deploy

`make deploy_grpc_health_check`

It will deploy deployment and service for tanjunchen/grpc-health-check:1.0.

## uninstall

`make delete_grpc_health_check`

It will delete deploy and service for tanjunchen/grpc-health-check:1.0.

## logs in pod/grpc-health-check

```
{
  "level": "info",
  "msg": "Serving the Check request for health check",
  "time": "2020-09-20T08:46:16Z"
}
```