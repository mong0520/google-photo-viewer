build:
	docker-compose build --no-cache

run_mongodb:
	docker-compose up -d mongo

run_redis:
	docker-compose up -d redis

run_al:
	docker-compose up