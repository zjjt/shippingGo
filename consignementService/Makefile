compile_proto:
	protoc  --go_out=. --micro_out=. proto/consignement/consignement.proto
build: 
		docker build -t shippinggo_consignementms .
run:
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 shippinggo_consignementms