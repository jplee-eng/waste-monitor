.PHONY: run build clean

run:
	go run cmd/server/main.go

build:
	go build -o waste-monitor cmd/server/main.go

clean:
	rm -f waste-monitor
	rm -f waste_monitor.db
