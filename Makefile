docker-build:
	docker build -t headless .

dev-setup:
    # Download kind from: http://kind.sigs.k8s.io
	kind create cluster --name headless-service
	kind load docker-image headless --name headless-service
	kubectl apply -f deploy/

run: docker-build dev-setup
	@echo "All done - you can use kubectl to explore the default namespace!"

teardown:
	kind delete cluster --name headless-service
