build: build_darwin_arm64 build_windows_386 build_windows_amd64 build_linux_amd64

build_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -o pkg/civ2lint_darwin_arm64

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -trimpath -o pkg/civ2lint_windows_amd64.exe

build_windows_386:
	GOOS=windows GOARCH=386 go build -trimpath -o pkg/civ2lint_windows_386.exe

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -trimpath -o pkg/civ2lint_linux_amd64
