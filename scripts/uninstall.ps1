# PowerShell Uninstall Script for fn - Fast Navigation Tool
param(
    [switch]$KeepData,
    [string]$InstallDir = "$env:LOCALAPPDATA\fn"
)

# Colors for output
function Write-Info { 
    param($Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warning { 
    param($Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error { 
    param($Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Get PowerShell profile info
function Get-PowerShellInfo {
    $profilePath = $PROFILE
    
    # Use CurrentUserAllHosts profile if available
    if ($PROFILE.CurrentUserAllHosts) {
        $profilePath = $PROFILE.CurrentUserAllHosts
    }
    
    return @{
        ProfilePath = $profilePath
        Shell = if ($PSVersionTable.PSEdition -eq "Core") { "pwsh" } else { "powershell" }
    }
}

# Remove binary and PATH entry
function Remove-Binary {
    param($InstallPath)
    
    $binaryPath = Join-Path $InstallPath "fn.exe"
    
    # Remove binary
    if (Test-Path $binaryPath) {
        Remove-Item $binaryPath -Force
        Write-Info "Removed binary: $binaryPath"
    }
    
    # Remove directory if empty
    if (Test-Path $InstallPath) {
        $items = Get-ChildItem $InstallPath
        if ($items.Count -eq 0) {
            Remove-Item $InstallPath -Force
            Write-Info "Removed directory: $InstallPath"
        }
    }
    
    # Remove from PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -like "*$InstallPath*") {
        $pathEntries = $currentPath -split ";" | Where-Object { $_ -ne $InstallPath -and $_ -ne "" }
        $newPath = $pathEntries -join ";"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Info "Removed from PATH: $InstallPath"
    }
}

# Remove PowerShell function from profile
function Remove-PowerShellFunction {
    $psInfo = Get-PowerShellInfo
    $profilePath = $psInfo.ProfilePath
    
    if (-not (Test-Path $profilePath)) {
        Write-Info "No PowerShell profile found"
        return
    }
    
    $content = Get-Content $profilePath -Raw
    
    if ($content -notmatch "function fn") {
        Write-Info "No fn function found in PowerShell profile"
        return
    }
    
    Write-Info "Removing fn function from PowerShell profile"
    
    # Remove the function and completion registration
    $lines = Get-Content $profilePath
    $newLines = @()
    $inFnFunction = $false
    $inCompletionBlock = $false
    
    foreach ($line in $lines) {
        if ($line -match "^# fn - Fast Navigation$") {
            $inFnFunction = $true
            continue
        }
        
        if ($inFnFunction) {
            if ($line -match "^function fn") {
                continue
            } elseif ($line -match "^Register-ArgumentCompleter.*fn") {
                $inCompletionBlock = $true
                continue
            } elseif ($inCompletionBlock -and $line -match "^\s*}") {
                $inCompletionBlock = $false
                continue
            } elseif ($inCompletionBlock) {
                continue
            } elseif ($line -match "^\s*}$" -and -not $inCompletionBlock) {
                $inFnFunction = $false
                continue
            } else {
                continue
            }
        }
        
        $newLines += $line
    }
    
    # Write back the modified content
    $newLines | Set-Content $profilePath -Encoding UTF8
    Write-Info "PowerShell function removed from profile"
}

# Remove data directory
function Remove-Data {
    $dataDir = Join-Path $env:USERPROFILE ".fn"
    
    if (Test-Path $dataDir) {
        if ($KeepData) {
            Write-Info "Keeping data directory: $dataDir"
        } else {
            Remove-Item $dataDir -Recurse -Force
            Write-Info "Removed data directory: $dataDir"
        }
    } else {
        Write-Info "No data directory found"
    }
}

# Main uninstallation
function Main {
    Write-Info "Uninstalling fn - Fast Navigation Tool"
    
    # Remove binary and PATH
    Remove-Binary $InstallDir
    
    # Remove PowerShell function
    Remove-PowerShellFunction
    
    # Remove data
    Remove-Data
    
    Write-Info "Uninstallation complete!"
    if (-not $KeepData) {
        Write-Info "All fn files and data have been removed"
    }
    Write-Info "Please restart PowerShell to complete the removal"
}

# Confirm uninstallation
$confirmation = Read-Host "Are you sure you want to uninstall fn? This will remove the binary, shell function, and$(if (-not $KeepData) { ' all bookmark data' }). (y/N)"
if ($confirmation -eq 'y' -or $confirmation -eq 'Y') {
    try {
        Main
    } catch {
        Write-Error "Uninstallation failed: $($_.Exception.Message)"
        exit 1
    }
} else {
    Write-Info "Uninstallation cancelled"
}