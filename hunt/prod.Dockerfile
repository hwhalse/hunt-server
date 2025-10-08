ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /src

COPY /services/hunt/go.mod /services/hunt/go.sum ./

RUN go mod download

COPY /services/hunt ./

RUN CGO_ENABLED=0 GOOS=linux go build -o hunt .

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /src/hunt /

USER nonroot:nonroot

CMD ["/hunt"]