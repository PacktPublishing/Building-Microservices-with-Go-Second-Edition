check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models


client-sdk: check-swagger
	cd client-sdk && GO111MODULE=off swagger generate client -f ../swagger.yaml