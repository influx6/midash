RUNENV := PORT=3000 API_VERSION=v1 MYSQL_PORT="5880" MYSQL_USER="dbUser" MYSQL_DATABASE=midash APP_DEV=dev MYSQL_PASSWORD="!thunder\#bot" 
ROOT= $(shell pwd)
RMD= $(ROOT)/../..

build:
	cd $(RMD)/ && make docker
	cd $(ROOT)
	docker build -t $(USER)/midash-dev ./

clean:
	docker rmi $(USER)/midash-dev

run: build
	docker run --rm --env-file vars.env --name $(USER)-midash -d $(USER)/midash-dev

stop:
	docker kill $(USER)-midash

shell: build
	docker run -it --rm --env-file vars.env $(USER)/midash-dev /bin/sh
