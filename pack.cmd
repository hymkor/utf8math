@echo off
set /P "VER=Version ? "

for %%I in (%CD%) do set "NAME=%%~nI"

for %%I in (linux windows.exe) do (
    for %%J in (386 amd64) do (
        set "GOOS=%%~nI"
        set "GOARCH=%%~J"
        go build -ldflags "-s -w"
        zip "%NAME%-%VER%-%%~nI-%%~J.zip" "%NAME%%%~xI"
    )
)
