export PATH := $(GOPATH)/bin:$(PATH)

version = 0.0.2
app = neotx

fmt:
	gofmt -w .

test:
	go test --cover ./...

build:
	rm -fr build 
	mkdir build 
	mkdir build/${app}-${version}-darwin 
	mkdir build/${app}-${version}-linux 
	mkdir build/${app}-${version}-win-386 
	mkdir build/${app}-${version}-win-amd64 
	env GOOS=darwin go build -o build/${app}-${version}-darwin/${app} .
	env GOOS=linux go build -o build/${app}-${version}-linux/${app} .
	env GOOS=windows GOARCH=amd64 go build -o build/${app}-${version}-win-amd64/${app}.exe . 
	env GOOS=windows GOARCH=386 go build -o build/${app}-${version}-win-386/${app}.exe .
	cp *.json build/${app}-${version}-darwin 
	cp *.json build/${app}-${version}-linux 
	cp *.json build/${app}-${version}-win-amd64 
	cp *.json build/${app}-${version}-win-386
	cd build && zip -rq ${app}-${version}-darwin.zip ${app}-${version}-darwin
	cd build && zip -rq ${app}-${version}-linux.zip ${app}-${version}-linux
	cd build && zip -rq ${app}-${version}-win-amd64.zip ${app}-${version}-win-amd64
	cd build && zip -rq ${app}-${version}-win-386.zip ${app}-${version}-win-386

.PHONY: build