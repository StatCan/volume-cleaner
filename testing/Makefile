apply:
	@minikube kubectl -- apply -f manifests.yaml
	@minikube kubectl -- apply -f ../manifests/scheduler/scheduler_config.yaml \
							-f  ../manifests/scheduler/scheduler_secret.yaml \
							-f  ../manifests/controller/controller_config.yaml \
							-f	  ../manifests/netpol.yaml \
							-f	  ../manifests/rbac.yaml \
							-f	  ../manifests/serviceaccount.yaml

sts:
	@minikube kubectl -- apply -f sts.yaml

delete_sts:
	@minikube kubectl -- delete sts sts1 -n anray-liu

ls:
	@minikube kubectl -- get sts -n anray-liu
	@minikube kubectl -- get pvc -n anray-liu

pvc1_logs:
	@minikube kubectl -- get -o yaml -n anray-liu pvc pvc1

pvc2_logs:
	@minikube kubectl -- get -o yaml -n anray-liu pvc pvc2

controller:
	@minikube kubectl -- apply -f controller/controller_deployment.yaml

stop_controller:
	@minikube kubectl -- delete -f controller/controller_deployment.yaml

controller_logs:
	@minikube kubectl -- logs -l app=volume-cleaner-controller -n das --tail 500

scheduler:
	@minikube kubectl -- apply -f scheduler/scheduler_job.yaml
	@minikube kubectl -- -n das wait --for=condition=complete job volume-cleaner-scheduler --timeout=300s || \
		(echo "Pod did not become ready"; exit 1)
	@make scheduler_logs
	
scheduler_logs:
	@minikube kubectl -- logs -l app=volume-cleaner-scheduler -n das

clean:
	@minikube kubectl delete ns anray-liu
	@minikube kubectl delete ns das

keep_ns:
	@minikube kubectl -- delete sts sts1 -n anray-liu
	@minikube kubectl -- delete pvc pvc1 -n anray-liu
	@minikube kubectl -- delete sts sts2 -n anray-liu
	@minikube kubectl -- delete pvc pvc2 -n anray-liu

das:
	@minikube kubectl -- -n das get pod

restart:
	@minikube delete
	@minikube start --driver=docker
	@make docker_build
	@make apply

restart_controller: clean apply controller


docker_build:
	@eval $(minikube docker-env) && \
	docker build -t controller -f ../docker/controller/Dockerfile .. && \
	docker build -t scheduler -f ../docker/scheduler/Dockerfile ..
	@minikube image load controller
	@minikube image load scheduler

docker_setup:
	@sudo groupadd docker 2>/dev/null || true
	@sudo usermod -aG docker lianra
	@newgrp docker
