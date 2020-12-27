sudo chmod 777 ../../code/res.go
./packfolder ../../code/view ../../code/res.go -v "res" -go
cd ../../code
go build -o ../app/simUI-64
cd ../app/
chmod u+x simUI-64
