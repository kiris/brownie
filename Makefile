GO111MODULES=on
APP?=stringifier
REGISTRY?=gcr.io/images

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

include development.mk
include docker.mk
