FROM node:latest as vuebuilder
WORKDIR /app
COPY files/private/vue ./
RUN npm install
COPY . .
RUN npm run build

FROM golang:1.14.0 as gobuilder
WORKDIR /go/src/app
COPY *.go ./
COPY files files/
RUN rm -r files/private/vue
COPY app app/
WORKDIR /go/src/app/app
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o app .

FROM scratch
ENV VERSION 1.0
MAINTAINER "Gian Marco Mennecozzi"
COPY --from=gobuilder /go/src/app/app/app .

EXPOSE 80
CMD ["./app"]