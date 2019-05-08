ARG VERSION=latest
FROM alpine:${VERSION} AS simplemicroservice
EXPOSE 10000/tcp
ARG VERSION
ARG SVC_NAME=microservice
ARG SVC_PATH=/service
RUN echo $VERSION > image_version
RUN mkdir ${SVC_PATH}
COPY ./${SVC_NAME} ${SVC_PATH}
# RUN chown nobody:nogroup ${SVC_PATH}/${SVC_NAME}
#CMD /service/microservice
