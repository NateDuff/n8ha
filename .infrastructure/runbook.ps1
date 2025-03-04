$ErrorActionPreference = "Stop"
$WarningPreference = "SilentlyContinue"

Connect-AzAccount -Identity -AccountId "b277e475-ba40-4952-849e-56d08df4443d" | Out-Null

$msiCert = Get-AzKeyVaultSecret -VaultName "kv-dcs-internal" -Name "MSICert" -AsPlainText

$certPath = New-TemporaryFile

Set-Content -Path $certPath.FullName -Force -AsByteStream -Value ([Convert]::FromBase64String($msiCert))

Disconnect-AzAccount

$certArgs = @($certPath.FullName, $null, [System.Security.Cryptography.X509Certificates.X509KeyStorageFlags]::MachineKeySet)

$cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2 -ArgumentList $certArgs

$connArgs = @{
    CertificateThumbprint = $cert.Thumbprint
    ApplicationId         = "a048a641-afb4-4666-831c-907a4f612331"
    TenantId              = "1d122f81-68ca-4aeb-91cd-0d1a6029e07b"
    Subscription          = "8c73818d-17aa-49c4-8876-c9a53f09ba11"
}

Connect-AzAccount -ServicePrincipal @connArgs | Out-Null

$tokenResults = Get-AzAccessToken -ResourceUrl "499b84ac-1321-427f-aa17-267ca6975798"

Write-Output $tokenResults.Token
