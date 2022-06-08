PACKAGE=consul-backup

bin-build:
	CGO_ENABLED=0 GOOS=linux go build -a -v -o  ${PACKAGE}  main.go

docker-build: bin-build
	docker build --tag consul-backup .
