gen-swagger:
	swag fmt && bash scripts/gen-swagger.sh

.PHONY:
	gen-swagger