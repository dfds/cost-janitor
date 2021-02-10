FROM node:current-buster as nodejs-builder
WORKDIR /frontend-build
COPY ./frontend/poc .

RUN npm install
RUN npm install -g yarn
RUN yarn run build


FROM golang:alpine as build

COPY . .
COPY --from=nodejs-builder /frontend-build/dist /frontend/poc/dist

WORKDIR /

RUN mkdir -p /app/dist && go build -i -o /app/dist/cost_janitor

FROM golang:alpine

COPY --from=build /app/dist/cost_janitor /app/cost_janitor

EXPOSE 8080

ENTRYPOINT ["/app/cost_janitor"]