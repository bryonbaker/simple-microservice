FROM golang as goimage
ENV SRC=/go/src/
RUN mkdir -p /go/src/
WORKDIR /go/src/go_docker
# RUN git clone -b master --single-branch https://github.com/bryonbaker/simple-microservice.git /go/src/go_docker/ \
RUN git clone https://github.com/bryonbaker/simple-microservice.git /go/src/go_docker/ \
# RUN git clone git@github.com:bryonbaker/simple-microservice.git /go/src/go_docker/ \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go get github.com/gorilla/mux
RUN go build -o bin/go_docker

ARG VERSION=latest
FROM alpine:${VERSION} AS microservice
ARG VERSION

# Put a container-image version identifier in the root directory.
RUN echo $VERSION > image_version

RUN apk add — no-cache bash
ENV WORK_DIR=/docker/bin
WORKDIR $WORK_DIR
COPY --from=goimage /go/src/go_docker/bin/ ./
ENTRYPOINT /docker/bin/go_docker
EXPOSE 10000