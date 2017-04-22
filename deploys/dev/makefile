RUNENV := PORT=3000 API_VERSION=v1 MYSQL_PORT="5880" MYSQL_USER="dbUser" MYSQL_DATABASE="midash" MYSQL_PASSWORD="!thunder#bot" APP_DEV=dev
ROOT= $(shell pwd)
RMD= $(ROOT)/../..

build:
	docker build -t $(USER)/midash-dev ./

clean:
	docker rmi $(USER)/midash-dev

boot:
	docker run --env-file vars.env -d $(USER)/midash-dev

local:
	build
	$(RUNENV) go run $(RMD)/cmd/midash/main.go