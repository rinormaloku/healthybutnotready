# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' -o /healthybutnotready

FROM gcr.io/distroless/static AS final
COPY --from=build --chown=nonroot:nonroot /healthybutnotready /healthybutnotready

EXPOSE 8080

CMD [ "/healthybutnotready" ]