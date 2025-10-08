ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /src

COPY . .

RUN go get . && \
    CGO_ENABLED=0 go build -a -installsuffix cgo \
	-ldflags "-s -w" \
	-o munc .

# Use distroless as minimal base image to package the binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /src/munc /

USER nonroot:nonroot

CMD ["/munc"]