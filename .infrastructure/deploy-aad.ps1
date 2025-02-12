## Script to run main.bicep to deploy rustdesk

param (
    $SubscriptionId = "435478a1-71b6-420b-a0af-73776f050e9d",
    $Location = "North Central US",
    [switch]$WhatIf
)

$ErrorActionPreference = "Stop"

Connect-AzAccount -Subscription $SubscriptionId -Tenant "c4afa44d-3294-4c39-b063-0593f7ae69d4" | Out-Null

$deployParams = @{
    Name = "deploy-home-aad"
    Location = $Location
    TemplateFile = "aad.bicep"
}

New-AzSubscriptionDeployment @deployParams -WhatIf

if ($WhatIf) {
    exit 0
}

New-AzSubscriptionDeployment @deployParams
