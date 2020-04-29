#build linux package
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build offlinesign_main.go;
mv offlinesign_main offlinesign_linux;

#build win64 package
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build offlinesign_main.go;
mv offlinesign_main.exe offlinesign_win64.exe;

#build linux mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build offlinesign_main.go;
mv offlinesign_main offlinesign_mac;