build:
	go build -v ./server.go

go:
	go vet ./...
	go build -v ./server.go
	./server


deploy:
	git rebase master deploy
	git status
	git push
	git switch master

.DEFAULT_GOAL := go