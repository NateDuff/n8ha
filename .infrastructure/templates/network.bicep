param location string
param environmentName string

resource publicIp 'Microsoft.Network/publicIPAddresses@2024-05-01' = {
  name: 'pip-${environmentName}'
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    publicIPAllocationMethod: 'Dynamic'
  }
}

resource vnet 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  name: 'vnet-${environmentName}'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.148.0.0/16'
      ]
    }
    subnets: [
      {
        name: 'containerAppsSubnet'
        properties: {
          addressPrefix: '10.148.0.0/24'
          delegations: [
            {
              name: 'delegation'
              properties: {
                serviceName: 'Microsoft.App/environments'
              }
            }
          ]
        }
      }
      {
        name: 'gatewaySubnet'
        properties: {
          addressPrefix: '10.148.255.0/27'
          
        }
      }
    ]
  }
}
