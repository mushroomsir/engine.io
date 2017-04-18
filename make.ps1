go test -v -coverprofile="transports.coverprofile" ./transports
go test -v -coverprofile="main.coverprofile"
D:\go\bin\gover
go tool cover -html="gover.coverprofile"
Remove-Item "*.coverprofile"