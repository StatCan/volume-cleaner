first: # prevents accidental running of make rules
	@echo "Please use explicit make commands with volume cleaner."

run-controller:
	@echo "🚧 Starting run..."
	@echo "🧰 Setting up run dependencies..."
	@kubectl apply -f manifests/rbac.yaml \
		-f manifests/serviceaccount.yaml \
		-f manifests/netpol.yaml \
		-f manifests/controller/controller_config.yaml
	@kubectl -n das apply -f manifests/controller/controller_deployment.yaml
	@echo "Ready to go!"

run-scheduler:
	@echo "🚧 Starting run..."
	@echo "🧰 Setting up run dependencies..."
	@kubectl apply -f manifests/rbac.yaml \
		-f manifests/serviceaccount.yaml \
		-f manifests/netpol.yaml \
		-f manifests/scheduler/scheduler_config.yaml
	@kubectl -n das apply -f manifests/scheduler/scheduler_job.yaml
	@kubectl -n das wait --for=condition=complete job volume-cleaner-scheduler --timeout=300s || \
		(echo "Pod did not become ready"; exit 1)
	@echo "Pod logs:"
	@kubectl -n das logs -l job-name=volume-cleaner-scheduler --tail 500

clean:
	@echo "🧼 Cleaning up leftover resources..."
	@kubectl delete -f manifests/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/controller/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/scheduler/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/controller/controller_deployment.yaml --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/scheduler/scheduler.job.yaml --ignore-not-found > /dev/null 2>&1 || true
	@echo "Cleaning complete"
