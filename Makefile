pwd = $(shell pwd)

.PHONY: db
db:
	docker container run --detach --name db \
			--restart always \
			--publish 3306:3306 \
			--env MYSQL_ROOT_PASSWORD=secret \
			--env MYSQL_DATABASE=snippetbox \
			--env MYSQL_USER=dev \
			--env MYSQL_PASSWORD=dev \
			--mount 'type=bind,src=$(pwd)/sql,dst=/docker-entrypoint-initdb.d' \
			mysql:8.0.28

.PHONY: run
run:
	go run ./cmd/web

.PHONY: testdb
testdb:
	docker container run --detach --name testdb \
			--restart always \
			--publish 3307:3306 \
			--env MYSQL_ROOT_PASSWORD=secret \
			--env MYSQL_DATABASE=snippetbox \
			--env MYSQL_USER=test \
			--env MYSQL_PASSWORD=test \
			mysql:8.0.28

.PHONY: test
test:
	go test -cover -v ./...
