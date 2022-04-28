FROM registry.access.redhat.com/ubi8/ubi:latest as build
WORKDIR /build

RUN dnf -y --disableplugin=subscription-manager install go

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o tenant-utils-fake . && strip tenant-utils-fake

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build /build/tenant-utils-fake /tenant-utils-fake
ENTRYPOINT ["/tenant-utils-fake"]
