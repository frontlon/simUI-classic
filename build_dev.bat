windres -o code/res.syso main.rc
cd code
go build -o ../app/simUI.exe
cd ../app/
simUI.exe

pause
