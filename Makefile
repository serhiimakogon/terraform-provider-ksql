NAME=ksql
BINARY=terraform-provider-${NAME}
VERSION=1.0.10-pre
OS_ARCH=darwin_arm64

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${NAME}/${VERSION}/${OS_ARCH}
