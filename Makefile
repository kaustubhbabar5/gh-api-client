run:
	docker-compose up

test:
	docker run -itd --rm -p 6379:6379 redis:7.0.2
	bash -c 'set -o allexport; source .env; set +o allexport; REDIS_PASSWORD="" go test ./... -v'