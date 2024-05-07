param(
    [Parameter(Position = 0)][string]$action,
    [Parameter(Position = 1)][string]$name,
    [Parameter(Position = 2)][string]$path

)

function Invoke-Add {
    param (
        [Parameter(Mandatory = $true)][string]$config,
        [Parameter(Mandatory = $true)][string]$name,
        [Parameter(Mandatory = $true)][string]$path
    )


    $config.envs += [PSCustomObject]@{
        name = $name
        path = $path
    }

    Write-Output ("Successfully added the new JDK: " + $name)
}

Import-Module $PSScriptRoot\add.ps1
Import-Module $PSScriptRoot\rm.ps1
Import-Module $PSScriptRoot\list.ps1
Import-Module $PSScriptRoot\use.ps1

if (!(Test-Path $PSScriptRoot\..\cache\config.json)) {
    New-Item -Path $PSScriptRoot -Name config.json -ItemType "file"
}

$config = Get-Content -Path $PSScriptRoot\..\cache\config.json -Raw | ConvertFrom-Json

if ($null -eq $config) {
    $config = New-Object -TypeName psobject
}

if (!($config | Get-Member envs)) {
    Add-Member -InputObject $config -MemberType NoteProperty -Name envs -Value @()
}


switch ($action) {
    add { Invoke-Add $config $name $path }
    rm { Invoke-Rm $config $name }
    list { Invoke-List $config }
    use { Invoke-Use $config $name }
    default {
        Write-Host 'jvm add <name> <path>   Add a JDK'
        Write-Host 'jvm rm <name>           Remove a JDK'
        Write-Host 'jvm list                List all installed JDKs'
        Write-Host 'jvm use <name>          Use a JDK' 
    }

}

ConvertTo-Json $config | Set-Content -Path $PSScriptRoot\..\cache\config.json


