.PHONY: build

build:
	sam build

dynamodb-local:
	docker compose up
dynamodb-local-create-table:
	aws dynamodb create-table --cli-input-json file://json/create-table.json --endpoint-url http://localhost:8000
dynamodb-local-add-items-work:
	aws dynamodb batch-write-item --cli-input-json file://json/work/add-items.json --endpoint-url http://localhost:8000
dynamodb-local-add-items-skills-tools:
	aws dynamodb batch-write-item --cli-input-json file://json/skills-tools/add-items.json --endpoint-url http://localhost:8000
start:
	sam.cmd local start-api --docker-network dynamodb-backend