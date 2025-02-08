$currentDir = Get-Location

Set-Location "$PSScriptRoot\.."

docker build -t duffraken -f Kraken.Dockerfile ./kraken

docker tag duffraken:latest duffnath/duffraken:latest

docker push duffnath/duffraken:latest

Set-Location $currentDir
