migrate:
	go-bindata -o ./pkg/postgres/migrations/bindata.go -pkg migrations pkg/postgres/migrations

run:
	go run cmd/upboat.go

test:
	go test ./... -covermode=count -v
