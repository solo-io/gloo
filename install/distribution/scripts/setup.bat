@echo off

SET SCRIPT_DIR=%~dp0
SET FILENAME="glooctl-windows-amd64.exe"


echo "        _                  _   _ "
echo "   __ _| | ___   ___   ___| |_| |"
echo "  / _\` | |/ _ \ / _ \ / __| __| |"
echo " | (_| | | (_) | (_) | (__| |_| |"
echo "  \__, |_|\___/ \___/ \___|\__|_|"
echo "  |___/                          "

:: Create .gloo directory in user home
mkdir %USERPROFILE%/.gloo/bin

:: Copy over glooctl Windows executable
xcopy %SCRIPT_DIR%/%FILENAME% %USERPROFILE%/.gloo/bin/

:: Echo instructions for next manual operation (we don't want to update the PATH here, as it's error prone on Windows)
echo Setup complete!
echo
echo Please add \"%USERPROFILE%/.gloo/bin/glooctl\" to your PATH environment variable