.PHONY: tidy-go-mod go-deps-update

tidy-go-mod:
	cd shared && go mod tidy

go-deps-update:
	cd shared && go get -u ./... && go mod tidy

