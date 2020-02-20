build:
	docker-compose -f ./docker-compose.yml build --parallel --force-rm
run:
	docker-compose -f ./docker-compose.yml up --detach --no-recreate
stop:
	docker-compose -f ./docker-compose.yml down --rmi all	
test:
	docker-compose -f ./docker-compose-test.yml up --build --abort-on-container-exit --exit-code-from banner_rotation_test
test-stop:
	docker-compose -f ./docker-compose-test.yml down --rmi all