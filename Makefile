.PHONY: all

docker-build:
	docker build -t tlantic/k8s-sidecar .

docker-run:
	docker run -d -p 50051:50051 -v ${KUBECONFIG_DIR}:/app/config --env KUBECONFIG=/app/config/kubeconfig --env K8S_NAMESPACE=${K8S_NAMESPACE} --name k8s-sidecar tlantic/k8s-sidecar