SOURCEDIR=cmd
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
MAIN := cmd/main.go

RELEASE_DIR=_output
BINARY=grisou

.DEFAULT_GOAL: build
.PHONY: build release clean

build:  $(SOURCES)
	go build -o $(BINARY) $(MAIN)

define build
	mkdir -o $(RELEASE_DIR)/$(1);
	GOOS=$(1) go build -o ${RELEASE_DIR}/$(1)/${BINARY} $(MAIN)
endef

release: $(SOURCES)
	$(call build, darwin)
	$(call build, linux)

clean:
	rm -f $(BINARY)
	rm -rf $(RELEASE_DIR)
