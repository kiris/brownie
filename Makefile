GO111MODULES=on
APP?=stringifier
REGISTRY?=gcr.io/images

ifneq (,$(wildcard .env))
	# include env vars from file
	include .env
	# export to makefile only env vars that are in file
	export $(shell sed 's/=.*//' .env)
endif


.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

#.env: .env.example
#	@echo 'run "cat .env.example > .env", then change the variables in the .env'
#	@exit 1

include makefiles/*.mk
