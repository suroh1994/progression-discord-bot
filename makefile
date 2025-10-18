

build-mbpg:
	cd Magic-Booster-Pack-Generator; \
	docker build \
	-t suroh/mbpg:latest \
	-f ./web.Dockerfile \
	.

run-mbpg:
	docker run \
	-d --rm \
	--name mbpg \
	-p 8080:8080 \
	suroh/mbpg:latest

stop-mbpg:
	docker kill mbpg

run-pgdb:
	docker run \
	-d --rm \
	--name card-db \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=progression \
	-p 5432:5432 \
	-v ./db-scripts:/docker-entrypoint-initdb.d \
	postgres:18.0-alpine3.22

stop-pgdb:
	docker kill card-db

run-dependencies: run-mbpg run-pgdb