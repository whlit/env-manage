function Invoke-Use {
    param (
        [Parameter(Mandatory = $true)][object]$config,
        [Parameter(Mandatory = $true)][string]$name
    )

    $env = $config.envs | Where-Object { $_.name -eq $name }

    if ($null -eq $env) {
        Write-Output ("JDK not found: " + $name)
        return
    } 
    [Environment]::SetEnvironmentVariable("JAVA_HOME", $env.path, 'User')

    Write-Output ("Successfully set JDK to: " + $name)

}