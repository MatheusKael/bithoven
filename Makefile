

build:
	cd ./cmd && go build  -o ../bin/bithoven main.go

run: build
	cd ./bin && ./bithoven
