@echo off

SET SCRIPT_DIR=%~dp0
SET FILENAME=glooctl-windows-amd64.exe


echo         _                  _   _
echo    __ _^| ^| ___   ___   ___^| ^|_^| ^|
echo   / _ `^| ^|/ _ \ / _ \ / __^| __^| ^|
echo  ^| (_^| ^| ^| (_) ^| (_) ^| (__^| ^|_^| ^|
echo   \__, ^|_^|\___/ \___/ \___^|\__^|_^|
echo   ^|___/
echo.

:: Create .gloo directory in user home
mkdir %USERPROFILE%\.gloo\bin
if %ERRORLEVEL% == 0 (
    echo Created directory %USERPROFILE%\.gloo
)
echo.


:: Copy over glooctl Windows executable
xcopy %SCRIPT_DIR%%FILENAME% %USERPROFILE%\.gloo\bin
if %ERRORLEVEL% == 0 (
    echo Copied glooctl executable to %USERPROFILE%\.gloo\bin
) else (
    echo Failed to copy glooctl executable to %USERPROFILE%\.gloo\bin
    goto :end
)

:: Echo instructions for next manual operation (we don't want to update the PATH here, as it's error prone on Windows)
echo.
echo Setup complete!
echo.
echo Please add the following directory to your PATH environment variable:
echo.
echo     %USERPROFILE%\.gloo\bin
echo.

:end