set GOOS=linux
set GOARCH=amd64
go build -o build/linux/amd64/go-hidden-service-bot ./
docker build . -t kvrg/go-tor-bot:latest

set GOOS=linux
set GOARCH=mips
set GOMIPS=softfloat
go build -o build/linux/mips/go-hidden-service-bot ./