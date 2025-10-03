BINARY=assignment-sensor
BUILD_DIR=build
SRC=./src
GO ?= $(shell command -v go)

.PHONY: all build install uninstall clean run fmt vet

all: build

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY) $(SRC)

install:  # assumes you've already run `make build`
	install -m 0755 $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	install -m 0644 systemd/assignment-sensor.service /etc/systemd/system/assignment-sensor.service
	systemctl daemon-reload

uninstall:
	- systemctl disable --now assignment-sensor.service
	 rm -f /usr/local/bin/$(BINARY)
	 rm -f /etc/systemd/system/assignment-sensor.service
	systemctl daemon-reload

run: build
	./$(BUILD_DIR)/$(BINARY) --interval=2s --device=internal

fmt:
	gofmt -w $(SRC)

vet:
	$(GO) vet $(SRC)

clean:
	rm -rf $(BUILD_DIR)
