build:
	go build -o myna -a 

docker:
	docker build -t myna .

clean:
	rm -rf myna

test:
	go test -v -cover ./...

