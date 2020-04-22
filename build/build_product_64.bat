packfolder.exe ../code/view ../code/res.go -v "res" -go
windres -o ../code/res.syso main.rc
cd ../code
go build -a -ldflags="-H windowsgui -w -s" -o ../app/simUI-64.exe
cd ../app/
simUI-64.exe
pause