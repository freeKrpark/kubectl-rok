.PHONY: build install clean

PLUGINS := kubectl-restart_history kubectl-images kubectl-drift kubectl-shell kubectl-clean

GOOS   := linux
GOARCH := amd64

build:
	@mkdir -p bin
	@for p in $(PLUGINS); do \
		echo "Building $$p..."; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$$p ./cmd/$$p || exit 1; \
	done

clean:
	rm -rf bin