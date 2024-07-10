all:run


run:
	go run server/cmd/main.go

	
mup: 
	go run  server/db/migrate/main.go up
mon:
	go run  server/db/migrate/main.go down

test:
	go test -timeout 30s -run ^TestE2E$$ backend/server/internal/user
	go test -timeout 30s -run ^TestE2E$$ backend/server/internal/product
	go test -timeout 30s -run ^TestE2E$$ backend/server/internal/cart
	go test -timeout 30s -run ^TestE2E$$ backend/server/internal/orders
