BINARY := tic-tac-claude

.PHONY: build run clean frontend

## build: install frontend deps, build frontend, then compile the Go binary
build: frontend
	go build -o $(BINARY) .

## frontend: install npm deps and compile the React app into dist/
frontend:
	cd frontend && npm install && npm run build

## run: full build then start the server
run: build
	./$(BINARY)

## clean: remove compiled outputs
clean:
	rm -rf dist $(BINARY) frontend/node_modules
