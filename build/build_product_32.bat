set GO111MODULE=on

z:
cd Z:\work\go\src\SimUI\build
windres -o ../code/res.syso main.rc
cd Z:\work\go\src\SimUI\code
go build -ldflags="-H windowsgui -w -s" -o ../app/simUI-32.exe
pause