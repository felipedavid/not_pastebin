postgres:
	docker run --name postgres -e "POSTGRES_PASSWORD=postgres" -p 5432:5432 -d postgres

build_image:
	docker build .

cert:
	mkdir tls && \
	cd tls && \
	go run /usr/lib/go/src/crypto/tls/generate_cert.go --host=localhost --rsa-bits=2048

.PHONY: postgres
