@echo off
setlocal enabledelayedexpansion

rem ---------- 参数 ----------
if "%~1"=="" (
  set "GO_VERSION=1.25.3"
) else (
  set "GO_VERSION=%~1"
)
if "%~2"=="" (
  set "REL_DEST=..\plugin\go\%GO_VERSION%"
) else (
  set "REL_DEST=%~2"
)

rem ---------- 环境与路径 ----------
rem 将脚本目录解析为绝对路径
set "SCRIPT_DIR=%~dp0"
rem 去掉末尾反斜杠（如果有）
if "%SCRIPT_DIR:~-1%"=="\" set "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"

rem 目标目录（相对于脚本目录）
for %%I in ("%SCRIPT_DIR%\%REL_DEST%") do set "DEST=%%~fI"

rem 构造下载 URL（windows-amd64）
set "ARCH=windows-amd64"
set "FILE=go%GO_VERSION%.%ARCH%.zip"
set "URL=https://go.dev/dl/%FILE%"

rem 临时下载文件
set "TMP_ZIP=%TEMP%\%FILE%"

echo.
echo Script dir: %SCRIPT_DIR%
echo Go version: %GO_VERSION%
echo Destination: %DEST%
echo Download URL: %URL%
echo Temp zip: %TMP_ZIP%
echo.

rem ---------- 创建目标目录 ----------
if not exist "%DEST%" (
  mkdir "%DEST%"
  if errorlevel 1 (
    echo Failed to create destination directory "%DEST%".
    exit /b 1
  )
)

rem ---------- 下载（使用 PowerShell 的 Invoke-WebRequest） ----------
echo Downloading %URL% ...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "try { Invoke-WebRequest -Uri '%URL%' -OutFile '%TMP_ZIP%' -UseBasicParsing -ErrorAction Stop; exit 0 } catch { Write-Error $_; exit 1 }"
if errorlevel 1 (
  echo Download failed. Please check network or URL.
  exit /b 1
)

rem ---------- 解压（使用 PowerShell Expand-Archive） ----------
echo Extracting to %DEST% ...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "try { if (Test-Path '%DEST%\go') { Remove-Item -Recurse -Force '%DEST%\go' } ; Expand-Archive -Path '%TMP_ZIP%' -DestinationPath '%DEST%' -Force; exit 0 } catch { Write-Error $_; exit 1 }"
if errorlevel 1 (
  echo Extraction failed.
  exit /b 1
)

rem ---------- 处理官方 zip 解压后包含的 go 子文件夹（将 go/* 移到 DEST 根） ----------
rem 有时 Expand-Archive 会在 DEST 中创建一个子目录 'go'，这里把其内容上移
if exist "%DEST%\go\" (
  echo Normalizing folder structure...
  powershell -NoProfile -ExecutionPolicy Bypass -Command ^
    "Get-ChildItem -Path '%DEST%\go' -Force | ForEach-Object { Move-Item -Path $_.FullName -Destination '%DEST%' -Force } ; Remove-Item -Recurse -Force '%DEST%\go'"
  if errorlevel 1 (
    echo Normalization failed.
    exit /b 1
  )
)

rem ---------- 清理临时 zip ----------
if exist "%TMP_ZIP%" del /q "%TMP_ZIP%"

rem ---------- 验证 go 可执行 ----------
set "GO_EXE=%DEST%\bin\go.exe"
if exist "%GO_EXE%" (
  echo.
  echo go executable found at: %GO_EXE%
  echo Running: "%GO_EXE%" version
  "%GO_EXE%" version
) else (
  echo go.exe not found under "%DEST%\bin\". Extraction may have failed.
  exit /b 1
)

rem ---------- 使用提示 ----------
echo.
echo 完成: Go %GO_VERSION% 已解压到:
echo   %DEST%
echo.
echo 临时会话使用（仅当前窗口有效）:
echo   set GOROOT=%DEST%
echo   set PATH=%DEST%\bin;%%PATH%%
echo.
echo 若要在脚本中立即使用该 go，请在调用命令时用绝对路径:
echo   "%DEST%\bin\go.exe" build -o myapp.exe ./cmd/myapp
echo.
echo 注意: 本脚本未修改系统环境变量, 若需永久生效请手动添加 GOROOT 和 PATH.
echo.

endlocal & pause
