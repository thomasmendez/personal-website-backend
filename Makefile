.PHONY: build

build:
	sam build

dynamodb-local:
	docker compose up
dynamodb-local-add-data:
	aws dynamodb create-table --cli-input-json file://json/create-table.json --endpoint-url http://localhost:8000
	aws dynamodb batch-write-item --cli-input-json file://json/add-table-items.json --endpoint-url http://localhost:8000

start:
	sam.cmd local start-api --docker-network dynamodb-backend