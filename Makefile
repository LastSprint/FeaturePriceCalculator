build:
	docker compose build --no-cache


start:
	docker compose up

rm:
	docker compose down
	docker compose rm