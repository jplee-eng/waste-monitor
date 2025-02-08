# waste-monitor

simple waste bin monitoring system using [LoRa](https://en.wikipedia.org/wiki/LoRa) and ultrasonic sensors

![interface screenshot](static/img/waste-monitor.jpeg)

## Overview

- monitors waste bin fill percentage, battery status, and signal strength
- realtime updates with [SSE](https://en.wikipedia.org/wiki/Server-sent_events)
- LoRa for longer range communication
- data persistence with sqlite
- tiny interface and only sqlite for dependencies

## Development

use the Makefile for easy dev

```sh
make run      # run development server
make build    # build binary
make test     # test out the sse and interface
make clean    # remove binary and db
```

### Test

to test the events sending to the frontend you can run `make build` in one terminal and `make test` from another and load up `http://localhost:8080` to see as the events send and update the site.