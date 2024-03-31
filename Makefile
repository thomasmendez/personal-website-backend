.PHONY: build

build:
	sam build

dynamo-local:
	docker compose up
	aws dynamodb create-table --cli-input-json file://json/create-table.json --endpoint-url http://localhost:8000
	aws dynamodb batch-write-item --cli-input-json file://json/add-table-items.json --endpoint-url http://localhost:8000

start:
	sam.cmd local start-api --docker-network dynamodb-backend