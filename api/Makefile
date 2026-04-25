
export-schema:
	docker compose run --rm -i -T sqldef -h database -W note -U note -p 5432 note --export | tee schema.sql

dry-run-schema:
	cat schema.sql | docker compose run --rm -i -T sqldef -h database -W note -U note -p 5432 note --dry-run

apply-schema:
	cat schema.sql | docker compose run --rm -i -T sqldef -h database -W note -U note -p 5432 note
