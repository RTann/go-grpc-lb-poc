all: images

# Protos
# ------

.PHONY: proto
proto: api/api.proto
	protoc \
	    --go_out=. --go_opt=paths=source_relative \
	    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	    $<

# Images
# ------

images-registry := quay.io/jvdm/stackrox-dev
images-tag-prefix := go-grpc-lb-poc

.PHONY: images
images: images-server images-client
images-%: images/%.Dockerfile
	podman build -f $< -t $(images-registry):$(images-tag-prefix)-$(<F:.Dockerfile=) .

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
