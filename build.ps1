# Build script for correlate CLI
# Usage: .\build.ps1 [version]
# Example: .\build.ps1 0.1.0-alpha

param(
    [string]$Version = "dev",
    [string]$Output = "correlate.exe"
)

# Get current git commit hash (if available)
$Commit = try { 
    git rev-parse --short HEAD 2>$null
} catch { 
    "none" 
}

# Get build date
$BuildDate = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

# Build with ldflags
$ldflags = "-X 'github.com/Rulopwd40/correlate/internal/commands.Version=$Version' " +
           "-X 'github.com/Rulopwd40/correlate/internal/commands.Commit=$Commit' " +
           "-X 'github.com/Rulopwd40/correlate/internal/commands.BuildDate=$BuildDate'"

Write-Host "Building correlate $Version..." -ForegroundColor Cyan
Write-Host "  Commit: $Commit"
Write-Host "  Date:   $BuildDate"
Write-Host ""

go build -ldflags $ldflags -o $Output ./cmd/correlate

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: $Output" -ForegroundColor Green
} else {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}
