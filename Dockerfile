FROM node:current-buster as nodejs-builder
WORKDIR /frontend-build
COPY ./frontend/poc .

RUN yarn install
RUN yarn run build


FROM golang:alpine as build
ARG HELLMAN_URL
ENV HELLMAN_URL=${HELLMAN_URL}

WORKDIR /src

COPY . .
COPY --from=nodejs-builder /frontend-build/dist /src/frontend/poc/dist

RUN go get github.com/GeertJohan/go.rice/rice
RUN cd src && rice embed-go
RUN mkdir -p /app/dist && cd src && go build -ldflags="-X 'main.HELLMAN_API_ENDPOINT=$HELLMAN_URL'" -i -o /app/dist/cost_janitor

FROM golang:alpine

COPY --from=build /app/dist/cost_janitor /app/cost_janitor

EXPOSE 8080

ENTRYPOINT ["/app/cost_janitor"]