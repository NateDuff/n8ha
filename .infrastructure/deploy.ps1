## Script to run main.bicep to deploy rustdesk

param (
    $SubscriptionId = "435478a1-71b6-420b-a0af-73776f050e9d",
    $Location = "North Central US",
    [switch]$WhatIf
)

$ErrorActionPreference = "Stop"

Connect-AzAccount -Subscription $SubscriptionId -UseDeviceAuthentication | Out-Null

$deployParams = @{
    Name = "deploy-rustdesk"
    Location = $Location
    TemplateFile = "main.bicep"
}

New-AzSubscriptionDeployment @deployParams -WhatIf

if ($WhatIf) {
    exit 0
}

New-AzSubscriptionDeployment @deployParams
