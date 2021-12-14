DOCKER		?= docker
PLATFORM 	?= linux/amd64
GOALS		:= image

.PHONY: $(GOALS)

image:
	$(DOCKER) build -t $(IMAGE_TAG) -f Dockerfile $(CONTEXT) --platform $(PLATFORM)
	$(DOCKER) push $(IMAGE_TAG)
