bindata:
	go-bindata -pkg migrations -ignore bindata -prefix ./datastore/migrations/ -o ./datastore/migrations/bindata.go ./datastore/migrations

table:
	migrate create -ext sql -dir ./datastore/migrations -seq create_merchants

createdb:
	psql -h localhost -U postgres -w -c "create database cuvva_merchant"

witchdb:
	psql -h localhost -U postgres -w -c "\c cuvva_merchant"