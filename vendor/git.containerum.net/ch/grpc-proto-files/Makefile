DOCKER_BUILDDIR = /src/ch-grpc
PROTOCFLAGS = -I/usr/include -I.
USE_DOCKER ?= 1
BUILDENV_CONTAINER_NAME = ch-grpc-buildenv
BUILDENV_DOCKERFILE = Dockerfile.buildenv
DOCKER_SOCKET = /var/run/docker.sock
DOCUMENTATION = DOCUMENTATION.md
DOCUMENTATION_FORMAT = markdown
GOPKG_PREFIX = git.containerum.net/ch/grpc-proto-files
SWAGGER_UI_DIR = swagger-ui

PACKAGES = common auth

sources = $(shell find * -name '*.proto')
have_docker_perms = $(shell docker images >/dev/null 2>/dev/null && echo 1)
run_prefix = 0
me = $(shell printf "%s:%s" $$(id -u) $$(id -g))
sources_in_package = $(shell find ${1} -name '*.proto')
outfiles_in_package = $(pathsubst %.proto,%${2},$(call sources_in_package,${1}))

.PHONY: clean go python doc swagger

all: go python swagger doc

## if USE_DOCKER is set to 1, run build inside docker container
ifeq (${USE_DOCKER},1)
$(info Using Docker container for building)
.dockerenv: $(BUILDENV_DOCKERFILE)
	@$(sudo) docker build -t ${BUILDENV_CONTAINER_NAME} -f $< . && touch $@
ifneq ($(have_docker_perms),1)
$(info You don`t have access to docker socket, so run with sudo)
sudo=sudo
endif
GOPATH=/src/go
run_prefix=$(sudo) docker run --rm \
		-v $(shell pwd):${DOCKER_BUILDDIR} \
		-w ${DOCKER_BUILDDIR} \
		${BUILDENV_CONTAINER_NAME}
endif

PROTOCFLAGS += -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

## build golang sources for gRPC
define build_go
	$(eval sources=$(call sources_in_package,${1}))
	$(run_prefix) protoc ${PROTOCFLAGS} --go_out=plugins=grpc:. $(sources)
	@$(sudo) chown $(me) -R ${GOPKG_PREFIX}/${1}
	@cp -rf ${GOPKG_PREFIX}/${1} .
	@$(run_prefix) find ${1} -name '*.pb.go' -exec protoc-go-inject-tag -input={} ';'

endef

define build_grpc_gateway
	$(eval sources=$(call sources_in_package,${1}))
	$(run_prefix) protoc ${PROTOCFLAGS} --grpc-gateway_out=. $(sources)
    @$(sudo) chown $(me) -R ${GOPKG_PREFIX}/${1}
    @cp -rf ${GOPKG_PREFIX}/${1} .
    @$(run_prefix) find ${1} -name '*.pb.gw.go' -exec protoc-go-inject-tag -input={} ';'

endef

go: .dockerenv
	$(foreach package,${PACKAGES},$(call build_go,$(package)))
	$(foreach package,${PACKAGES},$(call build_grpc_gateway,$(package)))

## build python sources for gRPC
define build_python
	$(eval sources=$(call sources_in_package,${1}))
	${run_prefix} python3 -m grpc.tools.protoc ${PROTOCFLAGS} --python_out=. --grpc_python_out=. $(call sources_in_package,${1})
	@$(sudo) chown $(me) -R $(sources:%.proto=%_pb2.py)
	@$(sudo) chown $(me) -R $(sources:%.proto=%_pb2_grpc.py)

endef

python: .dockerenv
	$(foreach package,${PACKAGES},$(call build_python,$(package)))

## generate swagger documentation
define build_swagger
	$(eval sources=$(call sources_in_package,${1}))
	@mkdir -p ${SWAGGER_UI_DIR}
	${run_prefix} protoc ${PROTOCFLAGS} --swagger_out=logtostderr=true:${SWAGGER_UI_DIR} $(sources)
	@$(sudo) chown $(me) -R $(sources:%.proto=${SWAGGER_UI_DIR}/%.swagger.json)

endef

swagger: .dockerenv
	$(foreach package,${PACKAGES},$(call build_swagger,$(package)))

## generate documentation
$(DOCUMENTATION): .dockerenv $(sources)
	@echo Generating documentation
	$(eval DOCUMENTATION=${DOCUMENTATION_FORMAT},${DOCUMENTATION})
	${run_prefix} protoc ${PROTOCFLAGS} --doc_out=. --doc_opt=${DOCUMENTATION} $(filter-out .dockerenv,$^)
	@$(sudo) chown $(me) -R $@

doc: $(DOCUMENTATION)

clean:
	@echo Cleaning...
	@rm -vf $(sources:.proto=_pb2.py)
	@rm -vf $(sources:.proto=_pb2_grpc.py)
	@rm -vf $(sources:.proto=.pb.go)
	@rm -vf $(sources:.proto=.pb.gw.go)
	@rm -vf $(sources:.proto=.swagger.json)
	@rm -rf ${SWAGGER_UI_DIR}
	@rm -vf ${DOCUMENTATION}

clean-env:
	@$(sudo) docker rmi -f ${BUILDENV_CONTAINER_NAME} || true
	@rm .dockerenv
