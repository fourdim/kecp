GO := go

.PHONY: build
build:
	CGO_ENABLED=0 $(GO) build -o ./build/kecp-server ./cmd/kecp-server
	mkdir -p ./build/app/dist
	pnpm -C ./app install
	pnpm -C ./app build
	cp -r ./app/dist/* ./build/app/dist/
	cp config.toml ./build/config.toml

.PHONY: pack
pack:
	tar -zcvf build/kecp.tar.gz -C build app/dist/* kecp-server

.PHONY: packc
packc:
	tar -zcvf build/kecp.tar.gz -C build config.toml app/dist/* kecp-server

.PHONY: run
run:
	cd build && ./kecp-server && cd ..
