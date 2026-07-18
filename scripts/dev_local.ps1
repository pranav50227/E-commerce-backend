# dev_local.ps1 - Starts the services locally in parallel using PowerShell background jobs.

$services = @(
    @{ Name = "UserManagementService"; Path = "service/UserManagementService"; Command = "go run cmd/server/main.go" },
    @{ Name = "ProductCatalogService"; Path = "service/ProductCatalogService"; Command = "go run cmd/server/main.go" },
    @{ Name = "InventoryService"; Path = "service/InventoryService"; Command = "go run cmd/server/main.go" },
    @{ Name = "OrderManagementService"; Path = "service/OrderManagementService"; Command = "go run cmd/server/main.go" },
    @{ Name = "ShoppingCartService"; Path = "service/ShoppingCartService"; Command = "go run cmd/server/main.go" },
    @{ Name = "auth-service"; Path = "service/auth-service"; Command = "go run cmd/server/main.go" },
    @{ Name = "api-gateway"; Path = "service/api-gateway"; Command = "go run cmd/gateway/main.go" }
)

$jobs = @()

foreach ($service in $services) {
    Write-Host "Starting $($service.Name) in background..."
    # Start-Job runs in a new session. We need to pass the environment if needed, or just let go run load .env if it does.
    # Actually, main.go loads .env relative to workspace. Let's make sure it does.
    $job = Start-Job -ScriptBlock {
        param($path)
        Set-Location $path
        go run (Get-ChildItem cmd -Recurse -Filter main.go | Select-Object -First 1).FullName
    } -ArgumentList (Resolve-Path $service.Path).Path
    $jobs += $job
}

Write-Host "------------------------------------------"
Write-Host "All services launched locally as PowerShell Jobs!"
Write-Host "To view jobs: Get-Job"
Write-Host "To view logs: Receive-Job -Job <Job>"
Write-Host "To stop services: Get-Job | Stop-Job"
Write-Host "------------------------------------------"
