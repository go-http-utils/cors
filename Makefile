test:
	go test -v -race

cover:
	rm -rf *.coverprofile
	go test -coverprofile=fresh.coverprofile
	gover
	go tool cover -html=fresh.coverprofile
	rm -rf *.coverprofile