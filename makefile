SHELL := /bin/bash

run:
	go run app/services/default-api/main.go | go run app/tooling/logfmt/main.go

update:
	go get -u -t -d -v ./...
	go mod vendor

VERSION := 1.0

all: default

default:
	docker build \
		-f zarf/docker/dockerfile \
		-t default-api-amd64:$(VERSION) \
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

KIND_CLUSTER := default

# Upgrade to latest Kind (>=v0.11): e.g. brew upgrade kind
# For full Kind v0.11 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.0
# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=default-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-default:
	kubectl get pods -o wide --watch --namespace=default-system

kind-load:
	cd zarf/k8s/kind/default-pod; kustomize edit set image default-api-image=default-api-amd64:$(VERSION)
	kind load docker-image default-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/default-pod | kubectl apply -f -

kind-restart:
	kubectl rollout restart deployment default-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-logs:
	kubectl logs -l app=default --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=default