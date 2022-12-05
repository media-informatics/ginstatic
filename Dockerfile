FROM golang:latest AS build

WORKDIR /ginstatic
COPY ./main.go .
RUN mkdir -p vendor
COPY go.mod .
COPY go.sum .
RUN go mod vendor
RUN go build -o StaticApp main.go

FROM debian
WORKDIR /app
COPY --from=build /ginstatic/StaticApp .
COPY static/ static/
COPY templates/ templates/
COPY seiten/ seiten/
EXPOSE 9000
CMD ["/app/StaticApp"]
