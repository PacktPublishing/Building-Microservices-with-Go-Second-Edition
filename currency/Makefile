.PHONY: protos

protos:
	 protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency
	 protoc -I proto_examples/ proto_examples/example.proto --go_out=plugins=grpc:proto_examples --go_opt=paths=source_relative