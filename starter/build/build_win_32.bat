windres -o ../res.syso main.rc
cd ..
go build -ldflags="-w -s" -o ./starter.exe
pause