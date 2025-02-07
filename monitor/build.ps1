$currentDir = Get-Location

Set-Location "$PSScriptRoot\.."

docker build -t pc-monitor . --build-arg FOLDERNAME=monitor --build-arg FINALIMAGE=gcr.io/distroless/static-debian12

docker tag pc-monitor:latest duffnath/duffhome-monitor:latest

docker push duffnath/duffhome-monitor:latest

Set-Location $currentDir
