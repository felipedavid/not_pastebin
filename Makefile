postgres:
	docker run --name postgres -e "POSTGRES_PASSWORD=postgres" -p 5432:5432 -d postgres

build_image:
	docker build .

.PHONY: postgres
