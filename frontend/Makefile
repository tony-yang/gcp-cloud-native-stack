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
	mkdir -p genproto
	protoc --go_out=plugins=grpc:genproto -I $(PROTODIR) $(PROTODIR)/app.proto

fmt:
	gofmt -s -w .

build: update-proto fmt
	go build -o /frontend .
