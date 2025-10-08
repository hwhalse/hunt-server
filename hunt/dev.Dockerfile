ARG GO_VERSION
FROM golang:${GO_VERSION} AS development

ENV PATH="/go/bin:$PATH"

WORKDIR /app

COPY services/hunt/go.mod services/hunt/go.sum services/hunt/.air.toml ./

RUN go install github.com/air-verse/air@latest && go mod download

COPY services/hunt .

CMD ["air", "-c", ".air.toml"]