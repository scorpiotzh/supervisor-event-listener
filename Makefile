# build file
GOCMD=go
# Use -a flag to prevent code cache problems.
GOBUILD=$(GOCMD) build -mod=vendor -ldflags -s -v -a

sup: BIN_BINARY_NAME=supervisor_listener
sup:
	GO111MODULE=on $(GOBUILD) -o $(BIN_BINARY_NAME) .

