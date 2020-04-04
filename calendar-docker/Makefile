run:
	docker-compose -f docker-compose.yml up --build
prune:
	docker-compose -f ./docker-compose.yml down --rmi local
prune-all:
	docker-compose -f ./docker-compose.yml down --rmi all

test:
	docker-compose -f docker-compose.test.yml up --build --exit-code-from calendar_test
test-prune:
	docker-compose -f docker-compose.test.yml down --rmi local
test-prune-all:
	docker-compose -f docker-compose.test.yml down --rmi all