all: build

clean:
	rm bin/praktikum-profile-setup

build: 
	go build -o bin/praktikum-profile-setup cmd/main.go
