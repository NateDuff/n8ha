targetScope = 'subscription'

@minLength(1)
@maxLength(64)
@description('Name of the environment that can be used as part of naming resource convention, the name of the resource group for your application will use this name, prefixed with rg-')
param environmentName string = 'n8rustdesk'

param location string = 'northcentralus'

resource rg 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: 'rg-${environmentName}'
  location: location
}

module rustDesk './rustDesk.bicep' = {
  name: 'rustDesk'
  scope: rg
  params: {
    location: location
    environmentName: environmentName
  }
}
