
NAME = infra-inventory

BUILD_DIR = ./build

clean:
	rm -f ${BUILD_DIR}
	
build:
	go get -d -v ./...
	go build -v -o ${BUILD_DIR}/${NAME} ./
	GOOS=linux go build -v -o ${BUILD_DIR}/${NAME}.linux ./

all: clean build
