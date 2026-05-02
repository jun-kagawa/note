SQLDEF_CMD := docker compose run --rm -i -T sqldef -h database -W note -U note -p 5432

export-schema:
	$(SQLDEF_CMD) --export | tee schema.sql

dry-run-schema:
	cat schema.sql | $(SQLDEF_CMD) note --dry-run

apply-schema:
	cat schema.sql | $(SQLDEF_CMD) note
	cat schema.sql | $(SQLDEF_CMD) note_test

setup:
	docker compose up -d | sleep 3

fmt:
	go fmt ./...
test: setup
	go test ./...
