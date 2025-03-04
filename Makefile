# Go related variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_VET=$(GO_CMD) vet
GO_CLEAN=$(GO_CMD) clean
GO_FMT=$(GO_CMD) fmt
GO_TEST=$(GO_CMD) test
LOCALSTACK_CMD=sh ./scripts/localstack.sh

start_collageapi:
	./bin/collageapi

run_collageapi:
	$(GO_CMD) run ./cmd/collageapi/main.go

run_migration: vet
	$(GO_CMD) run ./cmd/migrate/main.go

run_seed: vet
	$(GO_CMD) run ./cmd/seed/main.go

stack_deploy: 
	$(LOCALSTACK_CMD) -b

stack_clean:
	$(LOCALSTACK_CMD) -c

stack_start:
	$(LOCALSTACK_CMD) -s

collageapi: vet
	$(GO_BUILD) -C cmd/collageapi -o ../../bin/collageapi

vet: fmt
	$(GO_VET) ./...

fmt:
	$(GO_FMT) ./...

test:
	$(GO_TEST) --cover ./...

clean:
	$(GO_CLEAN) && \
	rm ./bin/* && \
	find ./resources -type f ! -name '.keep' -delete 


build: collageapi
start: vet build stack_deploy run_migration run_seed start_collageapi 
all: vet build

.PHONY: build fmt vet clean all

