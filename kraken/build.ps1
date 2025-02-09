param (
    [switch]$Preview
)

$tag = $Preview ? "preview" : "latest"

$currentDir = Get-Location

Set-Location "$PSScriptRoot\.."

docker build -t duffraken -f Kraken.Dockerfile ./kraken

docker tag duffraken:latest duffnath/duffraken:$tag

docker push duffnath/duffraken:$tag

Set-Location $currentDir
