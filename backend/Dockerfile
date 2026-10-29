# Development image only; not a production release artifact.
FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -trimpath -o /out/ploform-api ./cmd/ploform-api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/ploform-api /ploform-api
USER nonroot:nonroot
ENTRYPOINT ["/ploform-api"]
