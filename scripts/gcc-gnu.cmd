@echo off
setlocal
set "MINGW_BIN=C:\ProgramData\mingw64\mingw64\bin"
set "PATH=%MINGW_BIN%;%PATH%"
"C:\ProgramData\mingw64\mingw64\bin\x86_64-w64-mingw32-gcc.exe" %*
