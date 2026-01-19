gen-swagger:
	swag fmt && bash scripts/gen-swagger.sh
build: gen-swagger
	go build -ldflags="-w -s" -o out ./cmd/server

.PHONY:
	gen-swagger build