SOURCEDIR=cmd
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

RELEASE_DIR=_output
BINARY=grisou

.DEFAULT_GOAL: build
.PHONY: build release clean

build:
	go build -o $(BINARY) $(SOURCES)

define build
	mkdir -o $(RELEASE_DIR)/$(1);
	GOOS=$(1) go build -o ${RELEASE_DIR}/$(1)/${BINARY} $(SOURCES)
endef

release: $(SOURCES)
	$(call build, darwin)
	$(call build, linux)

clean:
	rm -f $(BINARY)
	rm -rf $(RELEASE_DIR)
