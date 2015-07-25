
NAME = infra-inventory

BUILD_DIR = ./build
BIN_DIR = ${BUILD_DIR}/usr/local/bin

clean:
	rm -rf ${BUILD_DIR}
	
build:
	go get -d -v ./...
	[ -d "${BIN_DIR}" ] || mkdir -p ${BIN_DIR}
	go build -v -o ${BIN_DIR}/${NAME}.$(shell go env GOOS) ./
	GOOS=linux go build -v -o ${BIN_DIR}/${NAME}.linux ./
	mkdir -p ${BUILD_DIR}/etc/${NAME}
	cp etc/${NAME}.json.sample ${BUILD_DIR}/etc/${NAME}/${NAME}.json

all: clean build
