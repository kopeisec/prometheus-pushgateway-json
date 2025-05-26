help:
	@cat Makefile | grep '# `' | grep -v '@cat Makefile'

# `make build`
.PHONY: build
build:
	docker build \
		--progress plain \
		--tag kopeisec/prometheus-pushgateway-json:latest .
