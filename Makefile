build: clean build-cli build-gui

build-cli: \
	pkg/civ2lint_darwin_arm64 \
	pkg/civ2lint_linux_amd64 \
	pkg/civ2lint_windows_386.exe \
	pkg/civ2lint_windows_amd64.exe

build-gui: \
	fyne-cross/dist/darwin-arm64/civ2lint.app \
	fyne-cross/dist/linux-amd64/civ2lint.tar.xz \
	fyne-cross/dist/windows-386/civ2lint.exe.zip \
	fyne-cross/dist/windows-amd64/civ2lint.exe.zip

fyne-cross/dist/darwin-arm64/civ2lint.app:
	fyne-cross darwin -arch arm64 -icon Icon.png -app-id com.github.leonsp.civ2lint

fyne-cross/dist/windows-amd64/civ2lint.exe.zip:
	fyne-cross windows -arch amd64 -icon Icon.png -app-id com.github.leonsp.civ2lint

fyne-cross/dist/windows-386/civ2lint.exe.zip:
	fyne-cross windows -arch 386 -icon Icon.png -app-id com.github.leonsp.civ2lint

fyne-cross/dist/linux-amd64/civ2lint.tar.xz:
	fyne-cross linux -arch amd64 -icon Icon.png

pkg/civ2lint_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -o pkg/civ2lint_darwin_arm64 cmd/civ2lint.go

pkg/civ2lint_windows_amd64.exe:
	GOOS=windows GOARCH=amd64 go build -trimpath -o pkg/civ2lint_windows_amd64.exe cmd/civ2lint.go

pkg/civ2lint_windows_386.exe:
	GOOS=windows GOARCH=386 go build -trimpath -o pkg/civ2lint_windows_386.exe cmd/civ2lint.go

pkg/civ2lint_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -trimpath -o pkg/civ2lint_linux_amd64 cmd/civ2lint.go

clean:
	rm -rfv pkg/* fyne-cross/dist/*

tools:
	go install fyne.io/tools/cmd/fyne@latest
	go install github.com/fyne-io/fyne-cross@latest
