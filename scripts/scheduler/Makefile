run_scheduler:
	@echo "🚧 Starting run..."
	@echo "🧰 Setting up run dependencies..."
	@kubectl apply -f ../../manifests/rbac.yaml \
		-f ../../manifests/serviceaccount.yaml \
		-f ../../manifests/netpol.yaml \
		-f ../../manifests/scheduler/scheduler_config.yaml \
		-f ../../manifests/scheduler/scheduler_secret.yaml
	@kubectl -n das apply -f ../../manifests/scheduler/scheduler_job.yaml
	@kubectl -n das wait --for=condition=complete job volume-cleaner-scheduler --timeout=300s || \
		(echo "Pod did not become ready"; exit 1)
	@echo "Pod logs:"
	@kubectl -n das logs -l job-name=volume-cleaner-scheduler --tail 500

