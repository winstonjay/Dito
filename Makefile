

all:
	go build -o dito -v
test:
	go test -v ./... --cover
clean:
	go clean
	rm -f dito
libTests:
	go run main.go lib/std_test.dito