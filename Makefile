.PHONY: install gen-mocks

install:
	go install .

gen-mocks:
	rm -frd mocks
	mockery --config configs/.mockery.yml
