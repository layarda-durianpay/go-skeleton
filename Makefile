# Get all keys from the YAML file as a space-separated list
KEYS := $(shell yq eval 'keys | join(" ")' application.yml)

OUTPUT = disbursement.exec


# Export each key dynamically
$(foreach key, $(KEYS), \
	$(eval export $(key)=$(shell yq eval '.$(key)' application.yml)))

build:
	go build -o $(OUTPUT) ./internal/cmd/main.go

start-service:
	make build && ./$(OUTPUT) start

start-consumer:
	make build && ./$(OUTPUT) consumer

clean:
	rm -f $OUTPUT

generate_openapi:
ifdef domainname
	@./scripts/openapi-http.sh ${domainname} internal/${domainname}/handler/http httphandler
else
	@echo "Please specify domainname, eg. 'make generate_openapi domainname=disburse'"
endif

generate_proto:
ifdef filename
	mkdir -p pkg/common/protogen

	protoc -I api/protobuf/ api/protobuf/${filename}.proto --go_out=:pkg/common/ --go-grpc_out=require_unimplemented_servers=false:pkg/common/ --experimental_allow_proto3_optional
else
	@echo "Please specify filename, eg. 'make generate_proto filename=disbursement'"
endif

db-create-migration:
ifdef filename
	migrate create -ext sql -dir database/sql_migrations -seq -digits 14 $(filename)
else
	@echo "Please specify filename, eg. 'make db-create-migration filename=add_column_name_in_table_name'"
endif

db-force-migrate:
ifdef version
	migrate -database $(POSTGRESQL_URL) -path database/sql_migrations force $(version)
else
	@echo "Please specify version, eg. 'make db-force-migrate version=000001'"
endif

db-curr-version:
	migrate -database $(POSTGRESQL_URL) -path database/sql_migrations version

db-migrate-down-level:
	migrate -database $(POSTGRESQL_URL) -path database/sql_migrations/ down $(filter-out $@,$(MAKECMDGOALS))

db_migrate_up:
	migrate -database $(POSTGRESQL_URL) -path database/sql_migrations/ up

db_migrate_down:
	migrate -database $(POSTGRESQL_URL) -path database/sql_migrations/ down