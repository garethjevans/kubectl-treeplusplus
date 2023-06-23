build:
	go build -o dist/kubectl-treeplusplus ./...

install: build
	cp dist/kubectl-treeplusplus /usr/local/bin
