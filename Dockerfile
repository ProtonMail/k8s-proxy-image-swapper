FROM golang:latest AS build

WORKDIR /
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build ./k8s-proxy-image-swapper .
ENTRYPOINT [ "./k8s-proxy-image-swapper" ]
