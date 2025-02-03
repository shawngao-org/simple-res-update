@ECHO OFF
timeout /t 3
move .\update-windows-amd64-tmp.exe .\update-windows-amd64.exe
.\update-windows-amd64.exe
