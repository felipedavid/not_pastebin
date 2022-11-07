createdb:
	 mariadb -u root -e "CREATE DATABASE not_pastebin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" -p

dropdb:
	mariadb -u root -e "DROP DATABASE not_pastebin;" -p

createuser:
	mariadb -u root -e "CREATE USER 'web'@'localhost';\
		GRANT SELECT, INSERT, UPDATE, DELETE ON not_pastebin.* TO 'web'@'localhost';\
		ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';" -p

dsn = "mysql://root:123@tcp(127.0.0.1:3306)/not_pastebin"

migrateup:
	migrate -path migration -database $(dsn) -verbose up

migratedown:
	migrate -path migration -database $(dsn) -verbose down

run:
	go run ./cmd/web

.PHONY: createdb dropdb createuser migrateup migratedown