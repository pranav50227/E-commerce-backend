# start_services.ps1 - Starts services locally in the background using persistent Windows processes.

$services = @(
    @{ Name = "UserManagementService"; Path = "service/UserManagementService"; Port = 8080 },
    @{ Name = "ProductCatalogService"; Path = "service/ProductCatalogService"; Port = 8081 },
    @{ Name = "InventoryService"; Path = "service/InventoryService"; Port = 8082 },
    @{ Name = "OrderManagementService"; Path = "service/OrderManagementService"; Port = 8083 },
    @{ Name = "ShoppingCartService"; Path = "service/ShoppingCartService"; Port = 8084 },
    @{ Name = "auth-service"; Path = "service/auth-service"; Port = 8085 },
    @{ Name = "api-gateway"; Path = "service/api-gateway"; Port = 8000 }
)

foreach ($service in $services) {
    $absPath = (Resolve-Path $service.Path).Path
    $port = $service.Port
    
    # Check if already running
    $conn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    if ($conn) {
        Write-Host "$($service.Name) is already running on port $port."
        continue
    }

    Write-Host "Starting $($service.Name) on port $port..."
    $mainFile = (Get-ChildItem -Path $absPath -Recurse -Filter main.go | Select-Object -First 1).FullName
    $outFile = Join-Path $absPath "stdout.log"
    $errFile = Join-Path $absPath "stderr.log"
    
    # Start the background process
    Start-Process -FilePath "go" -ArgumentList "run", "`"$mainFile`"" -WorkingDirectory $absPath -NoNewWindow -RedirectStandardOutput $outFile -RedirectStandardError $errFile
}

Write-Host "Waiting 5 seconds for services to initialize..."
Start-Sleep -Seconds 5

Write-Host "=== Verification ==="
$allOK = $true
foreach ($service in $services) {
    $conn = Get-NetTCPConnection -LocalPort $service.Port -ErrorAction SilentlyContinue
    if ($conn) {
        Write-Host "[OK] $($service.Name) is running on port $($service.Port)"
    } else {
        Write-Warning "[ERROR] $($service.Name) failed to start on port $($service.Port). Check logs: $($service.Path)/stdout.log and stderr.log"
        $allOK = $false
    }
}

if ($allOK) {
    Write-Host "All services started successfully!"
} else {
    Write-Error "Some services failed to start."
}
