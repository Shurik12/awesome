#!/bin/bash


build() {
	docker-compose build
}

start() {
	docker-compose up -d
}

stop() {
	docker-compose down
}

eval $1