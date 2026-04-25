SQLDEF_CMD := docker compose run --rm -i -T sqldef -h database -W note -U note -p 5432 note

export-schema:
	$(SQLDEF_CMD) --export | tee schema.sql

dry-run-schema:
	cat schema.sql | $(SQLDEF_CMD) --dry-run

apply-schema:
	cat schema.sql | $(SQLDEF_CMD)
