## Script to run main.bicep to deploy rustdesk

param (
    $SubscriptionId = "435478a1-71b6-420b-a0af-73776f050e9d",
    $Location = "North Central US",
    [switch]$WhatIf
)

$ErrorActionPreference = "Stop"

Connect-AzAccount -Subscription $SubscriptionId -Tenant "c4afa44d-3294-4c39-b063-0593f7ae69d4" | Out-Null

$deployParams = @{
    Name = "deploy-home-azures"
    Location = $Location
    TemplateFile = "main.bicep"
}

New-AzSubscriptionDeployment @deployParams -WhatIf

if ($WhatIf) {
    exit 0
}

New-AzSubscriptionDeployment @deployParams

# $subnetId = "/subscriptions/435478a1-71b6-420b-a0af-73776f050e9d/resourceGroups/rg-n8backupvm/providers/Microsoft.Network/virtualNetworks/vnet-n8backupvm/subnets/GatewaySubnet"
# $pipId = "/subscriptions/435478a1-71b6-420b-a0af-73776f050e9d/resourcegroups/rg-n8backupvm/providers/Microsoft.Network/publicIPAddresses/pip-n8backupvm"

# $gwipconfig = New-AzVirtualNetworkGatewayIpConfig -Name "vnet-gateway-config" -SubnetId $subnetId -PublicIpAddressId $pipId

# New-AzVirtualNetworkGateway -Name "vng-n8backupvm" -ResourceGroupName "rg-n8backupvm" -Location "North Central US" -IpConfigurations $gwipconfig -GatewayType "Vpn" -VpnType "RouteBased" -GatewaySku Basic
