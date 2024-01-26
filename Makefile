.PHONY: build run runfollower

build:
	go build -o bin/dcache

run: build
	./bin/dcache

runfollower: build
	./bin/dcache --port :4000 --leaderaddr :3000
