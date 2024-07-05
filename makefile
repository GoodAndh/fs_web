all:run


run:
	go run server/cmd/main.go

	
mup: 
	go run  server/db/migrate/main.go up
mon:
	go run  server/db/migrate/main.go down