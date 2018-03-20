gen:
	rm -R rpc
	mkdir -p rpc/blog
	myprotoc -I apidoc/proto apidoc/proto/*.proto --go_out=plugins=grpc:rpc/blog

run:
	go run main.go service article
