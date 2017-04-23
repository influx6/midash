RUNENV := PORT=3000 API_VERSION=v1 MYSQL_PORT="5880" MYSQL_USER="dbUser" MYSQL_DATABASE=midash APP_DEV=dev MYSQL_PASSWORD="!thunder\#bot" 
ROOT= $(shell pwd)
RMD= $(ROOT)/../..

build:
	cd $(RMD)/ && make
	cd $(ROOT)
	docker build -t $(USER)/midash-dev ./

clean:
	docker rmi $(USER)/midash-dev

run: build
	docker run --env-file vars.env -d $(USER)/midash-dev

shell: build
	docker run -it --rm --env-file  vars.env $(USER)/midash-dev /bin/sh
