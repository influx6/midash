build:
	docker build -t $(USER)/midash-dev ./

clean:
	docker rmi $(USER)/midash-dev

boot:
	docker run -d $(USER)/midash-dev