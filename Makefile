.PHONY: launch_services launch_services_with_tests stop_services build_services

launch_services: build_services
	docker compose up --force-recreate

launch_services_with_tests: build_services
	docker compose --profile test up --force-recreate --abort-on-container-exit --exit-code-from tester

stop_services:
	docker compose down -v

build_services:
	docker compose build
