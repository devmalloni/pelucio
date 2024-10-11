swag:
	swag init --parseDependency --parseInternal --parseDepth 100 -g cmd/http.go

go-client: 
	openapi-generator generate -i docs/swagger.yaml -g go -o pelucioclient -c docs/go.gen.yml