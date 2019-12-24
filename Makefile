test:	
	go test ./database ./runner ./utils ./cmd ./repositories/settings ./repositories/event ./read ./
cover:
	go test -coverprofile=coverage.out ./database ./runner ./utils ./cmd ./repositories/settings ./repositories/event ./read ./
	go tool cover -html=coverage.out -o coverage.html
coveralls:
	go test -v -covermode=count -coverprofile=coverage.out ./database ./runner ./utils ./cmd ./repositories/settings ./repositories/event ./read ./
	goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $(COVERALLS_TOKEN)
