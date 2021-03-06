FROM golang as builder

# !!! Docker layer caching will not repeat this step if the repo changes
# !!! You won't be able to build a test copy of your uncommitted code
RUN git clone --branch master https://github.com/bryonbaker/simple-microservice.git /go/src/go_docker
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers

# vvv Put magic environment variables in this line
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install go_docker
# ^^^

# Runtime image
FROM alpine:latest
COPY --from=builder /go/bin/go_docker /bin/go_docker
ARG VERSION=1.0
RUN echo $VERSION > /image_version
EXPOSE 10000
WORKDIR "/bin"
COPY ./boot-config.json /bin/
#COPY ./app-config.json /data/
CMD ["go_docker"]