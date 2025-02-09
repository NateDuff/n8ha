param (
    [switch]$Preview
)

$tag = $Preview ? "preview" : "latest"

$currentDir = Get-Location

Set-Location "$PSScriptRoot\.."

docker build -t pc-monitor -f Monitor.Dockerfile ./monitor

docker tag pc-monitor:latest duffnath/duffhome-monitor:$tag

docker push duffnath/duffhome-monitor:$tag

Set-Location $currentDir
