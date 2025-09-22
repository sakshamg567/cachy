export PATH := /usr/local/go/bin:$(PATH)

CACHE_PORTS = 50051 50052 50053
SERVER_PORT  = 8080
BIN_CACHE   = ./bin/cache-node
BIN_SERVER  = ./bin/server
PIDS_DIR    = .pids
LOGS_DIR    = logs

.PHONY: build-all run-all run-cache-all run-server stop-all clean logs

## Build binaries
build-all:
	go build -o $(BIN_CACHE) ./cmd/cache-node
	go build -o $(BIN_SERVER) ./cmd/server

## Run all servers
run-all: $(PIDS_DIR) $(LOGS_DIR)
	@> $(LOGS_DIR)/all.log   # truncate log
	$(MAKE) run-cache-all
	$(MAKE) run-server
	@echo "All servers running."

run-cache-all:
	@for port in $(CACHE_PORTS); do \
		echo "Starting cache-node on $$port"; \
		$(BIN_CACHE) --port $$port >> $(LOGS_DIR)/all.log 2>&1 & \
		sleep 0.1; \
		pid=$$(lsof -t -i :$$port | head -n1); \
		echo $$pid > $(PIDS_DIR)/cache-$$port.pid; \
	done


run-server:
	@echo "Starting Server on $(SERVER_PORT)"
	$(BIN_SERVER) --port $(SERVER_PORT) >> $(LOGS_DIR)/all.log 2>&1 & \
	pid=$$(lsof -t -i :$(SERVER_PORT) | head -n1); \
	echo $$pid > $(PIDS_DIR)/server.pid


## Stop all servers
stop-all:
	@echo "Stopping servers..."
	@if ls $(PIDS_DIR)/*.pid >/dev/null 2>&1; then \
		for pidfile in $(PIDS_DIR)/*.pid; do \
			if [ -f $$pidfile ]; then \
				name=$$(basename $$pidfile .pid); \
				kill -TERM `cat $$pidfile` 2>/dev/null && echo "Killed $$name (PID `cat $$pidfile`)"; \
				rm -f $$pidfile; \
			fi; \
		done; \
	fi
	@rmdir $(PIDS_DIR) 2>/dev/null || true

## Tail logs
tail-logs:
	tail -f logs/all.log | ccze -A

## Clean everything
clean: stop-all
	rm -rf $(PIDS_DIR) $(LOGS_DIR) ./bin

$(PIDS_DIR):
	mkdir -p $(PIDS_DIR)

$(LOGS_DIR):
	mkdir -p $(LOGS_DIR)
