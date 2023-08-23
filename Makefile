## Include variables from the .env file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: prin this help message
.PHONY: help
help: 
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ / '

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@DEV_PORT=${DEV_PORT} \
		PROD_PORT=${PROD_PORT} \
		DB_HOST=${DB_HOST} \
		DB_PORT=${DB_PORT} \
		DB_USER=${DB_USER} \
		DB_PASSWORD=${DB_PASSWORD} \
		DB_USER=${DB_USER} \
		DB_NAME=${DB_NAME} \
		JWT_SECRET=${JWT_SECRET} \
		CORS=${CORS} \
		go run .

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	# staticcheck ./...
	@echo 'Running tests'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building...'
	go build -ldflags='-s' -o=./bin/api .
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api .
