run-command:
	cd backend && PORT=8080 go run ./cmd/chat/command-api/main.go

#run-query:
#	cd backend && PORT=8081 go run ./cmd/chat/query-api/main.go

run-ws:
	cd backend && PORT=8082 go run ./cmd/chat/ws/main.go

run-backend:
	$(MAKE) run-command & $(MAKE) run-ws & wait
