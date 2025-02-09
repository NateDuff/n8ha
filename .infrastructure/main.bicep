targetScope = 'subscription'

param vpnEnvironmentName string = 'n8rustdesk'
param backupEnvironmentName string = 'n8backups'
param location string = 'northcentralus'

resource backupRg 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: 'rg-${backupEnvironmentName}'
  location: location
}

module backup './backup.bicep' = {
  name: 'backup'
  scope: backupRg
  params: {
    location: location
    environmentName: backupEnvironmentName
  }
}

// resource rustDeskRg 'Microsoft.Resources/resourceGroups@2022-09-01' = {
//   name: 'rg-${vpnEnvironmentName}'
//   location: location
// }

// module rustDesk './rustDesk.bicep' = {
//   name: 'rustDesk'
//   scope: rustDeskRg
//   params: {
//     location: location
//     environmentName: vpnEnvironmentName
//   }
// }
