.PHONY: all build build-plugins build-main clean

all: build

build: build-plugins build-main

build-plugins:
	@echo "Building plugins..."
	@mkdir -p plugins/yaml plugins/hcl
	go build -o plugins/yaml/plugin ./plugins/yaml/main.go
	go build -o plugins/hcl/plugin ./plugins/hcl/main.go

build-main:
	@echo "Building main application..."
	go build -o yaml2hcl2 main.go

clean:
	@echo "Cleaning up..."
	rm -f plugins/yaml/plugin
	rm -f plugins/hcl/plugin
	rm -f yaml2hcl2
	rm -f output.hcl output.json
