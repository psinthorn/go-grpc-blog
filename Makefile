
# auth-server start authentication server 
auth-server:
	go run ./backend/auth_service/service.go

# client start javascript frontend 
client:
	parcel ./frontend/index.html


# test run all test 
test:
	go test ./..

