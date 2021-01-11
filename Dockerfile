FROM golang:alpine as build

COPY src/ /src

WORKDIR /src

RUN mkdir -p /app/dist && go build -i -o /app/dist/cost_janitor

FROM golang:alpine

COPY --from=build /app/dist/cost_janitor /app/cost_janitor

EXPOSE 8080

ENTRYPOINT ["/app/cost_janitor"]