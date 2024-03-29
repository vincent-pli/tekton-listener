
# Image URL to use all building/pushing image targets
IMG ?= vincentpli/listener-controller:latest
IMG-LISTENER ?= vincentpli/listener:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

all: manager

go: manager docker-build docker-push deploy
# Run tests
test: generate fmt vet manifests
	go test ./api/... ./controllers/... -coverprofile cover.out

# Build manager binary
manager: generate fmt
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./api/...;./controllers/..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Build the docker image for listener
docker-build-listener:
	docker build . -f Dockerfile.listener -t ${IMG-LISTENER}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@value: .*@value: '"${IMG-LISTENER}"'@' ./config/default/manager_image_patch.yaml


# Push the docker image
docker-push:
	docker push ${IMG}
	docker push ${IMG-LISTENER}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-beta.1
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
