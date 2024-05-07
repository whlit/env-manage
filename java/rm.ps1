function Invoke-Rm {
    param (
        [Parameter(Mandatory = $true)][object]$config,
        [Parameter(Mandatory = $true)][string]$name
    )

    $config.envs = @($config.envs | Where-Object { $_.name -ne $name })

    Write-Output ("Successfully removed the JDK: " + $name)
}