FROM nginx:latest

ARG ARCH="amd64"
ARG OS="linux"

COPY ./docs/swagger-ui /usr/share/nginx/html

# docker build -f ./docs/Dockerfile -t functionx/fx-swagger-ui:1.0 . && docker push functionx/fx-swagger-ui:1.0