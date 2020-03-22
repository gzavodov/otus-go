run:
	docker-compose -f docker-compose.yml up --build
test:
	docker-compose -f docker-compose.test.yml up --build --exit-code-from calendar_test
stop:
	docker-compose -f ./docker-compose.yml down --rmi local
stop-all:
	docker-compose -f ./docker-compose.yml down --rmi all		