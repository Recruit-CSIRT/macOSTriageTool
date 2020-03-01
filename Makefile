
.PHONY: deps
deps:
	export GO111MODULE=on
	go mod -d
	go mod tidy

.PHONY: cli
cli:
	export GO111MODULE=on
	mkdir -p build
	go build -o ./build/mtt ./cmd/cli/

.PHONY: gui
gui:
	export QT_HOMEBREW=true
	export GO111MODULE=off
	qtdeploy build desktop ./cmd/gui/

	mkdir -p build
	rm -rf ./build/MTT.app

	cp -R ./automator/MTT.app ./build/
	cp -R ./cmd/gui/deploy/darwin/gui.app ./build/MTT.app/Contents/apps/

.PHONY: clean
clean:
	rm -rf ./cmd/gui/deploy
	rm -rf ./cmd/gui/darwin
	rm -rf ./build/*