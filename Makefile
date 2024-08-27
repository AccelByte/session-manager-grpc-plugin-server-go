SERVICE=session-manager-grpc-plugin-server
BUILDER_IMAGE=$(SERVICE)-builder
REVISION_ID=0.0.1
SERVICE_VERSION=0.0.1
export GIT_HASH?=unknown

RUN=docker run --rm \
	-v $(CURDIR):/opt/go/src/$(SERVICE) \
	-v $(GOPATH)/pkg/mod:/opt/go/pkg/mod \
	-w /opt/go/src/$(SERVICE)

proto:
	rm -rfv pkg/pb/*.pb.go
	mkdir -p pkg/pb
	docker run -t --rm -u $$(id -u):$$(id -g) -v $$(pwd):/data/ -w /data/ rvolosatovs/protoc:4.0.0 \
			--proto_path=pkg/proto --go_out=pkg/pb \
			--go_opt=paths=source_relative --go-grpc_out=pkg/pb \
			--go-grpc_opt=paths=source_relative ./pkg/proto/session-manager.proto

build:
	docker build -f Dockerfile.build -t $(BUILDER_IMAGE) .
	$(RUN) -e GOOS=linux $(BUILDER_IMAGE) \
		go build -buildvcs=false -o service -ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(date) -X main.gitHash=$(GIT_HASH) -X main.version=$(SERVICE_VERSION)" \
		./cmd/main.go
	docker build --no-cache --tag="$(SERVICE):$(REVISION_ID)" --tag="$(SERVICE):latest" .

run:
	docker run -e PLUGIN_GRPC_SERVER_AUTH_ENABLED=false -e $(AB_BASE_URL) -e $(AB_CLIENT_ID) -e $(AB_CLIENT_SECRET) $(SERVICE):latest