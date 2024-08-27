SERVICE=session-manager-grpc-plugin-server

proto:
	rm -rfv pkg/pb/*.pb.go
	mkdir -p pkg/pb
	docker run -t --rm -u $$(id -u):$$(id -g) -v $$(pwd):/data/ -w /data/ rvolosatovs/protoc:4.0.0 \
			--proto_path=pkg/proto --go_out=pkg/pb \
			--go_opt=paths=source_relative --go-grpc_out=pkg/pb \
			--go-grpc_opt=paths=source_relative ./pkg/proto/session-manager.proto

build-linux:
	docker build -f Dockerfile -t $(SERVICE) --build-arg="BUILDPLATFORM=linux/amd64" --build-arg="TARGETOS=linux" --build-arg="TARGETARCH=amd64" .

run:
	docker run -e PLUGIN_GRPC_SERVER_AUTH_ENABLED=false -e $(AB_BASE_URL) -e $(AB_CLIENT_ID) -e $(AB_CLIENT_SECRET) $(SERVICE):latest