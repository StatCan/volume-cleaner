first: # prevents accidental running of make rules
	@echo "Please use explicit make commands with volume cleaner."

dry-run: _dry-run-setup
	@echo "🚧 Starting dry run..."
	@kubectl -n das apply -f manifests/dry-run-job.yaml
	@echo "⏱️ Waiting for job to finish (up to 5 minutes)..."
	@kubectl -n das wait --for=condition=complete job/volume-cleaner-dry-run --timeout=300s || \
		(echo "❌ Pod did not become ready"; exit 1)
	@echo "📋 Pod logs:"
	@kubectl -n das logs -f -l job-name=volume-cleaner-dry-run
	@kubectl -n das delete -f manifests/dry-run-job.yaml || true
	@$(MAKE) clean
	@echo "✅ Dry run completed"

_dry-run-setup:
	@echo "🧰 Setting up dry-run dependencies..."
	@kubectl apply -f manifests/rbac.yaml \
		-f manifests/serviceaccount.yaml \
		-f manifests/netpol.yaml \
		-f manifests/dry-run-config.yaml

clean:
	@echo "🧼 Cleaning up leftover dry-run resources..."
	@kubectl delete -f manifests/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/dry-run-job.yaml --ignore-not-found > /dev/null 2>&1 || true
	@echo "Cleaning complete"
