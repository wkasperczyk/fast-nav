# PowerShell Install Script for fn - Fast Navigation Tool
param(
    [switch]$Force,
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

# Check if running as administrator for system-wide install
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Detect PowerShell version and profile
function Get-PowerShellInfo {
    $psVersion = $PSVersionTable.PSVersion.Major
    $profilePath = $PROFILE
    
    # Use CurrentUserAllHosts profile if available
    if ($PROFILE.CurrentUserAllHosts) {
        $profilePath = $PROFILE.CurrentUserAllHosts
    }
    
    return @{
        Version = $psVersion
        ProfilePath = $profilePath
        Shell = if ($PSVersionTable.PSEdition -eq "Core") { "pwsh" } else { "powershell" }
    }
}

# Install binary
function Install-Binary {
    param($InstallPath)
    
    # Check if binary exists
    if (-not (Test-Path "fast-nav.exe")) {
        Write-Error "Binary 'fast-nav.exe' not found. Please build it first."
        exit 1
    }
    
    # Create install directory
    if (-not (Test-Path $InstallPath)) {
        New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
        Write-Info "Created directory: $InstallPath"
    }
    
    # Copy binary
    $targetPath = Join-Path $InstallPath "fast-nav.exe"
    Copy-Item "fast-nav.exe" $targetPath -Force
    Write-Info "Binary installed: $targetPath"
    
    # Add to PATH if not already there
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$InstallPath*") {
        $newPath = "$currentPath;$InstallPath"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Info "Added to PATH: $InstallPath"
        Write-Warning "Please restart PowerShell to use the updated PATH"
    }
    
    return $targetPath
}

# Generate PowerShell function
function Get-PowerShellFunction {
    return @"

# fn - Fast Navigation
function fn {
    param([Parameter(ValueFromRemainingArguments)]`$args)
    
    if (`$args.Count -eq 0) {
        & fast-nav.exe
        return
    }
    
    `$firstArg = `$args[0]
    if (`$firstArg -in @('save', 'list', 'delete', 'path', 'edit', 'cleanup')) {
        & fast-nav.exe @args
    } else {
        `$dir = & fast-nav.exe navigate @args
        if (`$dir -and (Test-Path `$dir)) {
            Set-Location `$dir
        }
    }
}

# Tab completion for fn
Register-ArgumentCompleter -CommandName fn -ScriptBlock {
    param(`$commandName, `$parameterName, `$wordToComplete, `$commandAst, `$fakeBoundParameters)
    
    `$firstArg = `$commandAst.CommandElements[1].Value
    
    if (-not `$firstArg -or `$firstArg -eq `$wordToComplete) {
        # Complete main commands
        @('save', 'list', 'delete', 'path', 'edit', 'cleanup', 'navigate') | 
            Where-Object { `$_ -like "`$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new(`$_, `$_, 'ParameterValue', `$_) }
    } elseif (`$firstArg -in @('delete', 'path', 'navigate') -or (`$commandAst.CommandElements.Count -eq 2 -and `$firstArg -notin @('save', 'list', 'edit', 'cleanup'))) {
        # Complete aliases
        try {
            `$aliases = & fast-nav.exe list --quiet 2>`$null | ForEach-Object { `$_.Split()[0] }
            `$aliases | Where-Object { `$_ -like "`$wordToComplete*" } |
                ForEach-Object { [System.Management.Automation.CompletionResult]::new(`$_, `$_, 'ParameterValue', `$_) }
        } catch {
            # Ignore errors if fast-nav is not available or no bookmarks exist
        }
    }
}
"@
}

# Add PowerShell function to profile
function Add-PowerShellFunction {
    $psInfo = Get-PowerShellInfo
    $profilePath = $psInfo.ProfilePath
    
    Write-Info "Detected PowerShell $($psInfo.Version) ($($psInfo.Shell))"
    Write-Info "Profile: $profilePath"
    
    # Check if function already exists
    if (Test-Path $profilePath) {
        $content = Get-Content $profilePath -Raw
        if ($content -match "function fn") {
            if (-not $Force) {
                Write-Warning "PowerShell function already exists in profile"
                return
            } else {
                Write-Info "Replacing existing function (--Force specified)"
            }
        }
    }
    
    # Create profile directory if it doesn't exist
    $profileDir = Split-Path $profilePath -Parent
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
        Write-Info "Created profile directory: $profileDir"
    }
    
    # Add function to profile
    $function = Get-PowerShellFunction
    Add-Content -Path $profilePath -Value $function -Encoding UTF8
    
    Write-Info "PowerShell function added to profile"
    Write-Info "Please restart PowerShell or run: . `$PROFILE"
}

# Main installation
function Main {
    Write-Info "Installing fn - Fast Navigation Tool for Windows"
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-Error "PowerShell 5.0 or later is required"
        exit 1
    }
    
    # Install binary
    $binaryPath = Install-Binary $InstallDir
    
    # Add PowerShell function
    Add-PowerShellFunction
    
    Write-Info "Installation complete!"
    Write-Info ""
    Write-Info "Usage:"
    Write-Info "  fn save <alias>     - Save current directory"
    Write-Info "  fn <alias>          - Navigate to saved directory"
    Write-Info "  fn list             - List all bookmarks"
    Write-Info "  fn delete <alias>   - Delete a bookmark"
    Write-Info "  fn path <alias>     - Show path without navigating"
    Write-Info "  fn edit <alias>     - Update alias to current directory"
    Write-Info "  fn cleanup          - Remove dead bookmarks"
    Write-Info ""
    Write-Info "Note: Restart PowerShell to activate the function and updated PATH"
}

# Run main function
try {
    Main
} catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    exit 1
}