include make/proto.mk

all: images

# Protos
# ------

.PHONY: proto
proto: api/api.proto proto-install
	mkdir -p generated
	PATH=$(PROTO_BIN) && $(PROTOC) -I=api/ --go_out=generated/ --go-grpc_out=generated/ api/*

# Images
# ------

images-registry := quay.io/rtannenb/qa
images-tag-prefix := go-grpc-lb-poc

.PHONY: images
images: images-server images-client
images-%: images/%.Dockerfile
	docker build -f $< -t $(images-registry):$(images-tag-prefix)-$(<F:.Dockerfile=) .

# Deploy
# ------

deploy-specs :=  server-deployment.yaml \
		 server-service.yaml \
		 client-pod.yaml \
		 monitor.yaml

deploy-namespace := poc

.PHONY: deploy
deploy: $(addprefix specs/,$(deploy-specs))
	-oc create ns $(deploy-namespace)
	oc label --overwrite=true namespace/$(deploy-namespace) openshift.io/cluster-monitoring=true
	for f in $^; do kubectl -n $(deploy-namespace) apply -f $$f; done

.PHONY: teardown
teardown: $(addprefix specs/,$(deploy-specs))
	for f in $^; do kubectl -n $(deploy-namespace) delete -f $$f; done
	-oc delete ns $(deploy-namespace)
