IMAGE ?= "katalog-agent"
TAG ?= "dev"
MOCK_IMAGE ?= "katalog-server"

# Makefile
all: setup hooks

# requires `nvm use --lts` or `nvm use node`
.PHONY: setup
setup: 
	npm install -g @commitlint/config-conventional @commitlint/cli  

.PHONY: hooks
hooks:
	@git config --local core.hooksPath .githooks/

.PHONY: kind
kind:
	kind create cluster --name kind-katalog-agent --config kind/config.yaml

.PHONY: kind-delete
kind-delete:
	kind delete cluster --name kind-katalog-agent

.PHONY: kind-load
kind-load:
	docker build . -t $(IMAGE):$(TAG)
	kind load docker-image --name kind-katalog-agent $(IMAGE):$(TAG)
	cd docker/mockagentserver; docker build . -t $(MOCK_IMAGE):$(TAG)
	kind load docker-image --name kind-katalog-agent $(MOCK_IMAGE):$(TAG)


.PHONY: apply
apply:
	kustomize build k8s/ | kubectl apply -f -