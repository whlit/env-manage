function Invoke-Add {
    param (
        [Parameter(Mandatory = $true)][object]$config,
        [Parameter(Mandatory = $true)][string]$name,
        [Parameter(Mandatory = $true)][string]$path
    )

    foreach ($env in $config.envs) {
        if ($env.name -eq $name){
            Write-Output ("The name [" + $name + "] already exists. You can delete it and then add it again.   jvm rm <name>")
            return
        }
    }

    if (!(Test-Path $path)) {
        Write-Output ("The path [" + $path + "] does not exist.")
        return
    }

    if (!(Test-Path $path\bin\java.exe)) {
        Write-Output ("The path [" + $path + "] does not a java home directory.")
        return
    }

    $config.envs += @{
        name = $name
        path = $path
    }

    Write-Output ("Successfully added the new JDK: " + $name)
}