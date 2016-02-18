FROM golang:1.6

ENV GOPATH=/

RUN mkdir -p /src/portal-server/

COPY . /src/portal-server/

WORKDIR /src/portal-server/

RUN go build -o dbtool ./tool
RUN go build -o portalapi ./api
RUN go build -o portalgcm ./gcm

EXPOSE 8080
