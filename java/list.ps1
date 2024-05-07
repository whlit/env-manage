function Invoke-List {
    param (
        [Parameter(Mandatory = $true)][object]$config
    )


    Write-Host "All avaible versions of java"
    Write-Host ($config.envs | Format-Table | Out-String)
}