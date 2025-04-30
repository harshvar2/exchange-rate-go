BINARY_NAME = exchange-rate-service

download:
	go mod download
test-v: 
	go test -v -cover -covermode=atomic ./...
test:
	go test -cover -covermode=atomic ./...

list:
	go list ./... | grep -v /vendor/
lint:
	golangci-lint run -v
vet:
	go vet ./...
fmt:
	go fmt ./...
race:
	go test -race -short ./...
msan:
	go test -msan -short ./...

mocks:
	mockery --all

build:
	go build -o ${BINARY_NAME} cmd/server/main.go
docker-build:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose up --build --no-start
start-mysql: docker-build
	docker-compose start mysql
start-redis: docker-build
	docker-compose start cluster_initiator
start-app: docker-build
	docker-compose start app
restart-mysql:
	docker-compose restart mysql
restart-redis:
	docker-compose restart cluster_initiator
restart-app: docker-build
	docker-compose restart app
stop-app:
	docker-compose stop app
stop-redis:
	docker-compose stop cluster_initiator
stop-mysql:
	docker-compose stop mysql

live-reload:
	nodemon --exec go run cmd/server/main.go --signal SIGKILL --unhandled-rejections=strict SIGTERM
load-test:
	docker-compose -f locust-docker-compose.yml up locust-master

run: build
	./${BINARY_NAME}

test-run:
	go run cmd/server/main.go
sanity-test:
	k6 run tools/k6/test.js -e HOST=https://api.bobbleapp.asia 
k6-load-test:
	k6 run tools/k6/test.js -e HOST=https://api.bobbleapp.asia --vus 100 --duration 15m 
