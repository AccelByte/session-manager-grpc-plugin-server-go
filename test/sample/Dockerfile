FROM ngrok/ngrok as ngrok_image

FROM golang:1.22-alpine3.19

WORKDIR /workspace

RUN apk --no-cache add bash curl docker-cli-buildx git jq make

COPY --from=ngrok_image /bin/ngrok /bin/ngrok