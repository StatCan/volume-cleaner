run_controller:
	@echo "🚧 Starting run..."
	@echo "🧰 Setting up run dependencies..."
	@kubectl apply -f ../../manifests/rbac.yaml \
		-f ../../manifests/serviceaccount.yaml \
		-f ../../manifests/netpol.yaml \
		-f ../../manifests/controller/controller_config.yaml
	@kubectl -n zone apply -f ../../manifests/controller/controller_deployment.yaml
	@echo "Ready to go!"