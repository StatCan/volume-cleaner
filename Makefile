first: # prevents accidental running of make rules
	@echo "Please use explicit make commands with volume cleaner."

dry-run: _dry-run-setup
	@echo "🚧 Starting dry run..."
	@kubectl create namespace das || true
	@kubectl -n das apply -f dry-run-job.yaml
	@echo "⏱️ Waiting for job to start (up to 5 minutes)..."
	@kubectl -n das wait --for=condition=ready pod -l job-name=volume-cleaner-dry-run --timeout=300s || \
		(echo "❌ Pod did not become ready"; exit 1)
	@echo "📋 Pod logs:"
	@kubectl -n das logs -f -l job-name=volume-cleaner-dry-run
	@kubectl -n das delete -f dry-run-job.yaml || true
	@$(MAKE) stop
	@echo "✅ Dry run completed"

_dry-run-setup:
	@echo "🧰 Setting up dry-run dependencies..."
	@kubectl apply -f rbac.yaml \
		-f serviceaccount.yaml \
		-f netpol.yaml \
		-f dry-run-config.yaml

clean:
	@echo "🧼 Cleaning up leftover dry-run resources..."
	@kubectl delete -f rbac.yaml \
		-f serviceaccount.yaml \
		-f netpol.yaml \
		-f dry-run-config.yaml \
		-f dry-run-job.yaml \
		--ignore-not-found > /dev/null 2>&1 || true
