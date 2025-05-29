copy config.txt bin\config.txt
@REM Build without a cmd window opening.
go build -o bin\MyLocalhostGo.exe -ldflags -H=windowsgui .