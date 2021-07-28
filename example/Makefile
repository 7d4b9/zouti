
IMAGE := lever/example:$(shell git rev-parse --short -t HEAD)

export

image:
	@echo Building application image '$(IMAGE)'
	@docker build -t $(IMAGE) .
.PHONY: image

run-tests:
	@echo Launching functionnal tests
	@cd tests && go run .
.PHONY: run-tests