first: # prevents accidental running of make rules
	@echo "Please use explicit make commands with volume cleaner."

run:
	@echo "🚧 Starting run..."
	@echo "🧰 Setting up run dependencies..."
	@kubectl apply -f manifests/rbac.yaml \
		-f manifests/serviceaccount.yaml \
		-f manifests/netpol.yaml \
		-f manifests/controller/controller_config.yaml
	@kubectl -n das apply -f manifests/controller/controller_deployment.yaml
	@echo "Ready to go!"

clean:
	@echo "🧼 Cleaning up leftover resources..."
	@kubectl delete -f manifests/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/controller/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/controller/controller_deployment.yaml --ignore-not-found > /dev/null 2>&1 || true
	@echo "Cleaning complete"
