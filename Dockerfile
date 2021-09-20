FROM golang:1.17.1 AS build

WORKDIR /
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM scratch
USER 1000
COPY --from=build ./k8s-proxy-image-swapper .
ENTRYPOINT [ "./k8s-proxy-image-swapper" ]
