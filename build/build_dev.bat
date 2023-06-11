windres -o ../code/res.syso main.rc
cd ../code
D:\Program_Files\go\go1.20.4.windows-amd64\bin\go build -o ../app/simUI-64.exe
cd ../app/
simUI-64.exe

pause
