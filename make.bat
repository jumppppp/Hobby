
goversioninfo -config=./versioninfo.json -o=resource.syso
go build -o hobby.exe -buildmode=exe -ldflags="-s -w" -buildvcs=false -tags=netgo .\main.go