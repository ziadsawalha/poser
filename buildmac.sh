CGO_ENABLED=0 go build -a -ldflags '-s' server.go
mv server bin/macosx/poser
