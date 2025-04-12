CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -trimpath -buildvcs=false -o myapp
upx --lzma --best myapp