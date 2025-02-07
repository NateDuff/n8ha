$currentDir = Get-Location

Set-Location "$PSScriptRoot"

docker build -t pc-monitor .

docker tag pc-monitor:latest duffhome-monitor:latest

#docker push duffnath/duffhome-monitor:latest

Set-Location $currentDir
