.PHONY: build test clean release single_release

CMD_DIR:=cmd/mail-templater

# make directory and store path to variable
BUILDS_DIR:=$(PWD)/build
EXECUTABLE:=mail-templater
LDFLAGS=-X 'main.version=$(VERSION)' -w -s -extldflags '-static'

# go has build artifacts caching so soruce tracking not needed
build:
	@echo "Building mail-templater for current OS/architecture"
	@echo $(LDFLAGS)
	@CGO_ENABLED=0 go build -v -ldflags="$(LDFLAGS)" -o $(BUILDS_DIR)/$(EXECUTABLE) ./$(CMD_DIR)

build-for-docker:
	@echo $(LDFLAGS)
	@CGO_ENABLED=0 go build -v -ldflags="$(LDFLAGS)" -o  /tmp/$(EXECUTABLE) ./$(CMD_DIR)

test:
	@echo "Running tests"
	@go test -v ./...

clean:
	@rm -rf $(BUILDS_DIR)

install:
	@go install -ldflags="$(LDFLAGS)"

# lambda to generate build dir name using os,arch,version
temp_dir_name=$(EXECUTABLE)_$(1)_$(2)_v$(3)

pack_win=zip -r -j $(1).zip $(1) && rm -rf $(1)
pack_nix=tar -C $(1) -cpzf $(1).tar.gz ./ && rm -rf $(1)

define build_release
@echo "Building release package for OS $(1), arch $(2)"
$(eval temp_build_dir=$(BUILDS_DIR)/$(call temp_dir_name,$(1),$(2),$(VERSION)))
@mkdir -p $(temp_build_dir)
$(eval ifeq ($(1),windows)
	temp_executable=$(temp_build_dir)/$(EXECUTABLE).exe
else
	temp_executable=$(temp_build_dir)/$(EXECUTABLE)
endif)
@echo go build -tags="dev" -ldflags="$(RELEASE_LDFLAGS)"  -v -o $(temp_executable) ./$(CMD_DIR)
@GOOS=$(1) GOARCH=$(2) go build -tags="dev" -ldflags="$(RELEASE_LDFLAGS)" -v -o $(temp_executable) ./$(CMD_DIR)
$(eval ifeq ($(1),windows)
	pack_cmd = $(call pack_win,$(temp_build_dir))
else
	pack_cmd = $(call pack_nix,$(temp_build_dir))
endif)
@$(pack_cmd)
endef

release:
	$(call build_release,linux,amd64)
	$(call build_release,linux,386)
	$(call build_release,linux,arm)
	$(call build_release,darwin,amd64)
	$(call build_release,windows,amd64)
	$(call build_release,windows,386)

single_release:
	$(call build_release,$(OS),$(ARCH))

dev:
	@echo building $(VERSION)
	go build -v --tags="dev" --ldflags="$(LDFLAGS)" ./$(CMD_DIR)