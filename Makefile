COVERAGE=cover.out

$(COVERAGE):
	go test -v -coverprofile=$(COVERAGE) ./...

test: $(COVERAGE)

cover: test
	go tool cover -func=$(COVERAGE)

clean:
	rm -rf ./$(COVERAGE)

.PHONY: test cover clean