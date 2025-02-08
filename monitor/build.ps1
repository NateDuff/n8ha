$currentDir = Get-Location

Set-Location "$PSScriptRoot\.."

docker build -t pc-monitor -f Monitor.Dockerfile ./monitor

docker tag pc-monitor:latest duffnath/duffhome-monitor:latest

docker push duffnath/duffhome-monitor:latest

Set-Location $currentDir
