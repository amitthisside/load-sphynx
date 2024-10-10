@echo off
REM Loop through server names s1 to s5 and ports 5001 to 5005
setlocal enabledelayedexpansion

for %%i in (1 2 3 4 5) do (
    set "server_name=s%%i"
    set /a "port=5000 + %%i"
    echo Starting Python server: !server_name! on port !port!
    start cmd /k python server.py !server_name! !port!
)

echo All servers started.
pause
