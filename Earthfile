VERSION 0.6
FROM golang:1.17
ARG DOCKER_REPO=ghcr.io/dergeberl/kubeteach:dev
ARG BINPATH=/usr/local/bin/

deps:
    WORKDIR /src
    ENV GO111MODULE=on
    ENV CGO_ENABLED=0
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build-go:
    FROM +deps
    COPY --dir api/ pkg/ controllers/ main.go .
    ARG GOOS=linux
    ARG GOARCH=amd64
    ARG VARIANT
    RUN --mount=type=cache,target=/root/.cache/go-build \
        GOARM=${VARIANT#"v"} go build -ldflags="-w -s" -o build/kubeteach ./main.go
    SAVE ARTIFACT build/kubeteach

build-vue:
    FROM node:lts-alpine
    COPY --dir dashboard /dashboard
    WORKDIR /dashboard
    RUN npm install
    RUN npm run build
    SAVE ARTIFACT dist

docker:
    ARG TARGETPLATFORM
    ARG TARGETOS
    ARG TARGETVARIANT
    ARG TARGETARCH
    FROM --platform=$TARGETPLATFORM \
        gcr.io/distroless/static:nonroot
    #LABEL $DOCKER_LABEL
    # use the following to not for multiarch with emulation as desribed in
    # https://docs.earthly.dev/docs/guides/multi-platform#creating-multi-platform-images-without-emulation
    COPY --platform=linux/amd64 \
        (+build-go/kubeteach --GOOS=$TARGETOS --GOARCH=$TARGETARCH --VARIANT=$TARGETVARIANT) /kubeteach
    COPY --platform=linux/amd64 (+build-vue/dist) /dashboard
    USER 65532:65532
    ENTRYPOINT ["/kubeteach"]
    SAVE IMAGE --push $DOCKER_REPO

multiarch-docker:
    BUILD --platform=linux/amd64 +docker
    BUILD --platform=linux/arm/v7 +docker
    BUILD --platform=linux/arm64 +docker

test:
    ARG KUBERNETES_VERSION=1.23.x
    FROM +deps
    COPY +gotools/bin/setup-envtest $BINPATH
    RUN setup-envtest use $KUBERNETES_VERSION
    COPY --dir api/ pkg/ controllers/ crds/ main.go .
    COPY +build-vue/dist dashboard/dist
    #ARG GO_TEST="go test -race -coverprofile cover.out ./..."
    RUN eval `setup-envtest use -p env $KUBERNETES_VERSION` && \
        CGO_ENABLED=1 go test -race -coverprofile cover.out ./...
    SAVE ARTIFACT cover.out AS LOCAL cover.out

coverage:
    FROM +deps
    COPY --dir api/ pkg/ controllers/ crds/ main.go .
    COPY +test/cover.out .
    RUN go tool cover -func=cover.out

lint:
    FROM +deps
    COPY +gotools/bin/golangci-lint $BINPATH
    COPY --dir api/ pkg/ controllers/ crds/ main.go .golangci.yaml .
    RUN golangci-lint run -v ./...

manifests:
    FROM +deps
    COPY +gotools/bin/* $BINPATH
    COPY --dir api/ pkg/ controllers/ main.go .
    RUN controller-gen crd paths="./..." output:crd:artifacts:config=crds
    SAVE ARTIFACT crds AS LOCAL crds

check-manifests:
    FROM +deps
    COPY --dir crds/ .
    COPY +manifests/crds crds-new
    RUN diff crds crds-new

generate:
    FROM +deps
    COPY +gotools/bin/* $BINPATH
    COPY --dir api/ pkg/ controllers/ hack/ main.go .
    RUN controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
    SAVE ARTIFACT api AS LOCAL api

check-generate:
    FROM +deps
    COPY --dir api/ .
    COPY +generate/api api-new
    RUN diff -r api api-new


all:
    BUILD +deps
    BUILD +generate
    BUILD +manifests
    BUILD +lint
    BUILD +coverage
    BUILD +multiarch-docker

gotools:
    RUN GOBIN=/go/bin go get sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
    RUN GOBIN=/go/bin go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
    RUN GOBIN=/go/bin go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0
    SAVE ARTIFACT /go/bin

ci:
    BUILD +deps
    BUILD +check-generate
    BUILD +check-manifests
    BUILD +lint
    BUILD +test
    BUILD +coverage