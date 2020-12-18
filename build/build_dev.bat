windres -o ../code/res.syso main.rc
cd ../code
go build -o ../app/simUI-64.exe
cd ../app/
simUI-64.exe -env rBtHsZ

pause
