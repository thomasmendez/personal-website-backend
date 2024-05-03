.PHONY: build

db:
	docker compose up
db-create-table:
	aws dynamodb create-table --cli-input-json file://json/create-table.json --endpoint-url http://localhost:8000
db-add-items-work:
	aws dynamodb batch-write-item --cli-input-json file://json/work/add-items.json --endpoint-url http://localhost:8000
db-add-items-skills-tools:
	aws dynamodb batch-write-item --cli-input-json file://json/skills-tools/add-items.json --endpoint-url http://localhost:8000
build:
	sam build
start:
	sam.cmd local start-api --docker-network dynamodb-backend