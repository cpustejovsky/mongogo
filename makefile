SHELL := /bin/bash

run:
	go run app/services/starter-api/main.go | go run app/tooling/logfmt/main.go

update:
	go get -u -t -d -v ./...
	go mod vendor

VERSION := 1.0

all: mongogo

mongogo:
	docker build \
		-f zarf/docker/dockerfile \
		-t mongogo-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := mongogo

# Upgrade to latest Kind (>=v0.11): e.g. brew upgrade kind
# For full Kind v0.11 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.0
# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=mongogo-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-mongogo:
	kubectl get pods -o wide --watch --namespace=mongogo-system

kind-load:
	kind load docker-image mongogo-amd64:$(VERSION) --name $(KIND_CLUSTER)
	# cd zarf/k8s/kind/mongogo-pod; kustomize edit set image mongogo-api-image=mongogo-api-amd64:$(VERSION)

kind-apply:
	kustomize build zarf/k8s/kind/mongogo-pod | kubectl apply -f -

kind-restart:
	kubectl rollout restart deployment mongogo-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-logs:
	kubectl logs -l app=mongogo --all-containers=true -f --tail=100
	# kubectl logs -l app=mongogo --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=mongogo