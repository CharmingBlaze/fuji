@echo off
setlocal
set "MINGW_BIN=C:\ProgramData\mingw64\mingw64\bin"
set "PATH=%MINGW_BIN%;%PATH%"
"C:\Program Files\LLVM\bin\clang.exe" --target=x86_64-w64-windows-gnu %*
