first: # prevents accidental running of make rules
	@echo "Please use explicit make commands with volume cleaner."

run_controller:
	@make -C scripts/controller run_controller

run_scheduler:
	@make -C scripts/scheduler run_scheduler

create_job:
	@kubectl create job volume-cleaner-scheduler --from=cronjob/volume-cleaner-scheduler -n das

test:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	open cover.html

clean:
	@echo "🧼 Cleaning up leftover resources..."
	@kubectl delete -f manifests/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/controller/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/scheduler/ --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/controller/controller_deployment.yaml --ignore-not-found > /dev/null 2>&1 || true
	@kubectl delete -f manifests/scheduler/scheduler.job.yaml --ignore-not-found > /dev/null 2>&1 || true
	@echo "Cleaning complete"
