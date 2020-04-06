compile_proto:
	protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/vessel/vessel.proto
build: 
		docker build -t shippinggo_vesselms .
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 shippinggo_vesselms