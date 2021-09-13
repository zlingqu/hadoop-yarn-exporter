BINARY_NAME=hadoop-yarn-exporter
tag=`date "+%Y-%m-%d-%H-%M"`
adpBackendImage=docker.dm-ai.cn/cc/hadoop-yarn-exporter:$(tag)


clean:
	rm -rf build
compile: clean init-env
	cd $(srcDir) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) -v ./

run:
	cd build && ./$(BINARY_NAME)

docker-build: compile
	docker build -t $(adpBackendImage) -f service-run.Dockerfile .

docker-push:
	docker push $(adpBackendImage)