@echo off
chcp 936 >nul
setlocal enabledelayedexpansion

rem ---------- 参数与默认值 ----------
if "%~1"=="" (
  set "GOROOT_REL=..\plugin\go\1.25.3"
) else (
  set "GOROOT_REL=%~1"
)
if "%~2"=="" (
  set "OUTDIR=..\plugin\wails3"
) else (
  set "OUTDIR=%~2"
)
if "%~3"=="" (
  set "WAILS_REF=main"
) else (
  set "WAILS_REF=%~3"
)

rem ---------- 解析绝对路径 ----------
for %%I in ("%~dp0%GOROOT_REL%") do set "GOROOT=%%~fI"
for %%I in ("%~dp0%OUTDIR%") do set "GOBIN=%%~fI"

echo Using portable Go: "%GOROOT%\bin\go.exe"
echo Output directory: "%GOBIN%"
echo Wails ref: %WAILS_REF%
echo.

rem ---------- 前置检查 ----------
if not exist "%GOROOT%\bin\go.exe" (
  echo ERROR: go.exe not found in "%GOROOT%\bin\".
  echo Please download and extract Go to the GOROOT path or provide correct GOROOT_REL.
  pause
  exit /b 2
)

if not exist "%GOBIN%" (
  mkdir "%GOBIN%" 2>nul
  if errorlevel 1 (
    echo ERROR: failed to create output directory "%GOBIN%".
    pause
    exit /b 3
  )
)

rem ---------- 临时环境设置 ----------
set "OLD_PATH=%PATH%"
set "PATH=%GOROOT%\bin;%PATH%"
set "GOBIN=%GOBIN%"
set "GOROOT=%GOROOT%"
rem 可根据需要启用 GOPROXY
set "GOPROXY=https://goproxy.cn,direct"

rem ---------- 尝试用 go install 安装（优先） ----------
echo Trying go install (v3 path)...
"%GOROOT%\bin\go.exe" install github.com/wailsapp/wails/v3/cmd/wails3@latest
if %ERRORLEVEL%==0 goto :installed

echo go install (v3 path) failed, trying non-v3 path...
"%GOROOT%\bin\go.exe" install github.com/wailsapp/wails/cmd/wails3@latest
if %ERRORLEVEL%==0 goto :installed

echo go install attempts failed. Will fallback to clone + build.
goto :clone_build

:installed
rem 安装成功，检查可执行位置
echo go install succeeded. Searching for wails executable...
if exist "%GOBIN%\wails.exe" (
  set "WAILS_EXE=%GOBIN%\wails.exe"
) else (
  rem 检查默认 GOPATH\bin
  for /f "usebackq delims=" %%P in (`"%GOROOT%\bin\go.exe" env GOPATH`) do set "GOPATH_VAL=%%P"
  if exist "%GOBIN%\wails.exe" set "WAILS_EXE=%GOBIN%\wails.exe"
  if not defined WAILS_EXE if exist "%GOPATH_VAL%\bin\wails.exe" set "WAILS_EXE=%GOPATH_VAL%\bin\wails.exe"
  if not defined WAILS_EXE if exist "%GOROOT%\bin\wails.exe" set "WAILS_EXE=%GOROOT%\bin\wails.exe"
)

if defined WAILS_EXE (
  echo Found wails executable: "%WAILS_EXE%"
  echo Version:
  "%WAILS_EXE%" version
  goto :done
) else (
  echo go install reported success but wails.exe not found in common locations.
  goto :clone_build
)

:clone_build
rem ---------- 回退：克隆源码并用局部 go build ----------
echo Cloning wails repository and building from source...
set "SRC_DIR=%~dp0wails-src"

if not exist "%SRC_DIR%\.git" (
  echo Cloning to "%SRC_DIR%"...
  git clone https://github.com/wailsapp/wails.git "%SRC_DIR%"
  if errorlevel 1 (
    echo ERROR: git clone failed. Ensure git is installed and network is available.
    goto :restore_env
  )
) else (
  echo Repo exists, fetching updates...
  pushd "%SRC_DIR%"
  git fetch --all --tags
  git reset --hard
  git clean -fd
  popd
)

pushd "%SRC_DIR%"
git checkout %WAILS_REF% 2>nul || git checkout main
git pull origin %WAILS_REF% 2>nul || echo Note: pull failed or ref not present, building current checkout.

rem 确保 cmd/wails 存在
if not exist ".\cmd\wails\main.go" (
  echo ERROR: cmd/wails not found in source tree. Check repo structure or ref.
  popd
  goto :restore_env
)

echo Building wails CLI to "%GOBIN%\wails.exe" ...
"%GOROOT%\bin\go.exe" build -o "%GOBIN%\wails.exe" ./cmd/wails
if errorlevel 1 (
  echo ERROR: go build failed. Check build errors above.
  popd
  goto :restore_env
)

popd

if exist "%GOBIN%\wails.exe" (
  echo Build succeeded: "%GOBIN%\wails.exe"
  "%GOBIN%\wails.exe" version
  goto :done
) else (
  echo ERROR: wails.exe not found after build.
  goto :restore_env
)

:restore_env
rem 恢复 PATH
set "PATH=%OLD_PATH%"
echo Environment restored.
pause
exit /b 4

:done
rem 恢复 PATH 并结束
set "PATH=%OLD_PATH%"
echo.
echo Completed. wails.exe is available at:
for %%F in ("%GOBIN%\wails.exe") do if exist "%%~fF" echo   %%