param location string
param environmentName string
param fileShareName string = 'backups'

resource storageAccount 'Microsoft.Storage/storageAccounts@2023-01-01' = {
#disable-next-line BCP334
  name: 'st${environmentName}'
  location: location
  kind: 'StorageV2'
  sku: { name: 'Standard_LRS' }

  resource fileServices 'fileServices' = {
    name: 'default'

    resource fileShare 'shares' = {
      name: fileShareName
      properties: {}
    }
  }
}
