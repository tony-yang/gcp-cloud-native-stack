PROTODIR=../proto

test:
	go test -v

cover:
	go test -coverprofile=coverage.out -v
	@echo "########################################################"
	@echo "# Coverage Report                                      #"
	@echo "########################################################"
	go tool cover -func=coverage.out

update-proto:
	protoc --go_out=plugins=grpc:genproto -I $(PROTODIR) $(PROTODIR)/app.proto

build: update-proto
	go build -o /catalog .
