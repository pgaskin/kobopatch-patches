@echo off
cd "%~dp0"
erase "KoboRoot.tgz" >nul 2>&1
"./bin/koboptch-windows.exe"