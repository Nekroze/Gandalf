#!/bin/sh
set -euf

build() {
	docker-compose build --force-rm
}

mmock() {
	docker-compose up --remove-orphans -d mock
	#docker-compose up mockgen
}

run() {
	docker-compose up --force-recreate --remove-orphans --exit-code-from gandalf gandalf
}

down() {
	docker-compose rm -svf
	docker-compose down --rmi local --volumes --remove-orphans
}

cleanup() {
	exitcode=$?
	if [ "$exitcode" -ne 0 ]; then
		docker-compose logs mock
		echo Press enter when you are ready to clean all artifacts
		read -r
	fi
	find . -name '*.json' -delete || true
	down
	export COMPOSE_FILE=
	export COMPOSE_PROJECT_NAME=
}

setup() {
	export COMPOSE_FILE=examples/prototype/docker-compose.yml
	export COMPOSE_PROJECT_NAME=gandalf
	trap cleanup EXIT
}

setup && build && mmock && run
