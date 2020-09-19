.PHONY: deploy_grpc_health_check
deploy_grpc_health_check:
	@echo "... Building docker image ..."
	docker build . -t tanjunchen/grpc-health-check:1.0
	@echo "... Deploying in kubernetes ..."
	kubectl apply -f ./deploy/grpc-health-check.yaml

.PHONY: delete_grpc_health_check
delete_grpc_health_check:
	@echo "... Deleting tanjunchen-grpc-health-check-service ..."
	kubectl delete svc tanjunchen-grpc-health-check
	@echo "... Deleting tanjunchen-grpc-health-check-deployment ..."
	kubectl delete deployment tanjunchen-grpc-health-check