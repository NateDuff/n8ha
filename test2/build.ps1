param (
    [switch]$Preview
)

$tag = $Preview ? "preview" : "latest"

$currentDir = Get-Location

Set-Location "$PSScriptRoot\.."

docker build -t dufftest -f Test.Dockerfile ./test2

docker tag dufftest:latest duffnath/hello-web:$tag

docker push duffnath/hello-web:$tag

Set-Location $currentDir
