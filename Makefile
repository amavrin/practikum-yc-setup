all: build

clean:
	rm bin/praktikum-profile-setup

build: 
	CGO_ENABLED=0 go build -o bin/praktikum-profile-setup cmd/main.go
