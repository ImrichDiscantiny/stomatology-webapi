# use specific versions of images
FROM openapitools/openapi-generator-cli:v7.0.1 as api

WORKDIR /local

COPY api api
COPY scripts scripts
COPY .openapi-generator-ignore .openapi-generator-ignore

RUN docker-entrypoint.sh generate -c /local/scripts/generator-cfg.yaml

# not used normally but redefine entrypoint for the case of checking this stage results
ENTRYPOINT ["bash"]

############################################
FROM golang:latest AS build

WORKDIR /app

# download dependencies - low frequency of changes
COPY go.mod go.sum ./
RUN go mod download

# copy sources - higher frequency of changes
COPY internal/ internal/
COPY cmd/ cmd/
COPY --from=api /local/ ./


# create executable - stomatology-webapi-srv
# we want to use scratch image so setting
# the build options in the way that will link all dependencies statically
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
      -ldflags="-w -s" \
      -installsuffix 'static' \
      -o ./stomatology-webapi-srv ./cmd/stomatology-api-service

############################################
FROM scratch

# see https://github.com/opencontainers/image-spec/blob/main/annotations.md for details
LABEL org.opencontainers.image.authors="ID"
LABEL org.opencontainers.image.title="Stomatology Appointment List WebAPI Service"
LABEL org.opencontainers.image.description="WEBAPI for managing entries in ambulances` waiting list"

# list all variables and their default values for clarity
ENV STOMATOLOGY_API_ENVIRONMENT=production
ENV STOMATOLOGY_API_PORT=8080
ENV STOMATOLOGY_API_MONGODB_HOST=mongo
ENV STOMATOLOGY_API_MONGODB_PORT=27017
ENV STOMATOLOGY_MONGODB_DATABASE=id-stomatology
ENV STOMATOLOGY_API_MONGODB_COLLECTION=stomatology
ENV STOMATOLOGY_API_MONGODB_USERNAME=root
ENV STOMATOLOGY_API_MONGODB_PASSWORD=
ENV STOMATOLOGY_API_MONGODB_TIMEOUT_SECONDS=5

COPY --from=build /app/stomatology-webapi-srv ./

# Actual port may be changed during runtime
# Default using for the simple case scenario
EXPOSE 8080
ENTRYPOINT ["./stomatology-webapi-srv"]