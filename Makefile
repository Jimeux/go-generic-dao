.PHONY: test
test:
	go test -p 1 ./...

.PHONY: db-init
db-init:
	docker exec -i ${DATABASE_CONTAINER} mysql -u${DATABASE_USER} -p${DATABASE_PASSWORD} < db/init.sql
	docker exec -i ${DATABASE_CONTAINER} mysql -u${DATABASE_USER} -p${DATABASE_PASSWORD} ${DATABASE_NAME} < db/schema.sql
	docker exec -i ${DATABASE_CONTAINER} mysql -u${DATABASE_USER} -p${DATABASE_PASSWORD} ${DATABASE_NAME}test < db/schema.sql
