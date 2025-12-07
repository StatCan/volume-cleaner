run_scheduler:
	@echo "🚧 Starting run..."
	@echo "🧰 Setting up run dependencies..."
	@kubectl apply -f ../../manifests/rbac.yaml \
		-f ../../manifests/serviceaccount.yaml \
		-f ../../manifests/netpol.yaml \
		-f ../../manifests/scheduler/scheduler_config.yaml \
		-f ../../manifests/scheduler/scheduler_secret.yaml
	@kubectl -n zone apply -f ../../manifests/scheduler/scheduler_job.yaml
	@echo "Ready to go!"
