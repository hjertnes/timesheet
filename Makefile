test:
	go test -coverprofile=coverage.out ./database ./runner ./utils ./cmd ./repositories/settings ./repositories/event ./read ./
	go tool cover -html=coverage.out -o coverage.html
