IMAGE := zouti/example:$(shell git rev-parse --short -t HEAD)

export

example-image:
	$(info Building example application image '$(IMAGE)')
	@docker build -t $(IMAGE) -f example.Dockerfile .
.PHONY: example-image

example-run-tests:
	$(info Launching functionnal tests)
	@cd example/tests && go run .
.PHONY: example-run-tests