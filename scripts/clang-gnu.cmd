@echo off
setlocal
set "MINGW_BIN=C:\ProgramData\mingw64\mingw64\bin"
set "LLVM_BIN=C:\Program Files\LLVM\bin"
if exist "%LLVM_BIN%\lld.exe" (
  set "PATH=%LLVM_BIN%;%MINGW_BIN%;%PATH%"
) else (
  set "PATH=%MINGW_BIN%;%PATH%"
)
"C:\Program Files\LLVM\bin\clang.exe" --target=x86_64-w64-windows-gnu -fuse-ld=lld -Wno-override-module %*
