FROM golang as goimage
ENV SRC=/go/src/
RUN mkdir -p /go/src/
WORKDIR /go/src/go_docker
# RUN git clone -b master --single-branch https://github.com/bryonbaker/simple-microservice.git /go/src/go_docker/ \
RUN git clone https://github.com/bryonbaker/simple-microservice.git /go/src/go_docker/ \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go get github.com/gorilla/mux
RUN go build -o bin/go_docker

FROM alpine:latest AS microservice
RUN apk --no-cache add bash
ENV WORK_DIR=/docker/bin
WORKDIR $WORK_DIR
COPY --from=goimage /go/src/go_docker/bin/ ./

# Put a container-image version identifier in the root directory.
ARG VERSION=1.0
RUN echo $VERSION > image_version

EXPOSE 10000
#CMD ["go_docker"]