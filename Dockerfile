ARG BUILDPLATFORM
FROM --platform=$BUILDPLATFORM rvolosatovs/protoc:4.0.0 as proto
WORKDIR /build
COPY pkg/proto pkg/proto
RUN mkdir -p pkg/pb
RUN protoc --proto_path=pkg/proto --go_out=pkg/pb \
            --go_opt=paths=source_relative --go-grpc_out=pkg/pb \
            --go-grpc_opt=paths=source_relative pkg/proto/session-manager.proto

FROM --platform=$BUILDPLATFORM golang:1.23.0-alpine as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=proto /build/pkg/pb pkg/pb
RUN env GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o session-manager-grpc-plugin-go_$TARGETOS-$TARGETARCH cmd/main.go

FROM alpine:3.20
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY --from=builder /build/session-manager-grpc-plugin-go_$TARGETOS-$TARGETARCH session-manager-grpc-plugin-go
# Plugin arch gRPC server port
EXPOSE 6565
# expose port 8080 later this for metric
EXPOSE 8080
CMD [ "/app/session-manager-grpc-plugin-go" ]
