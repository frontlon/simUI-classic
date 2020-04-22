packfolder.exe ../code/view ../code/res.go -v "res" -go
windres -o ../code/res.syso main.rc
cd ../code
go build -ldflags="-H windowsgui -w -s" -o ../app/simUI-32.exe
cd ../app/
simUI-32.exe
pause