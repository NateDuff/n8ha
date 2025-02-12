param location string
param environmentName string
param fileShareName string = 'rustdeskdata'

resource storageAccount 'Microsoft.Storage/storageAccounts@2023-01-01' = {
#disable-next-line BCP334
  name: 'sg${environmentName}'
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

resource vnet 'Microsoft.Network/virtualNetworks@2024-05-01' = {
  name: 'vnet-${environmentName}'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        '10.0.0.0/16'
      ]
    }
    subnets: [
      {
        name: 'containerAppsSubnet'
        properties: {
          addressPrefix: '10.0.0.0/23'
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
    ]
  }
}
resource containerAppEnv 'Microsoft.App/managedEnvironments@2023-05-01' = {
  name: 'cae-${environmentName}'
  location: location
  properties: {
    workloadProfiles: [{
      workloadProfileType: 'Consumption'
      name: 'consumption'
    }]
    vnetConfiguration: {
      infrastructureSubnetId: vnet.properties.subnets[0].id
    }
  }

  resource storage 'storages' = {
    name: 'azurefiles'
    properties: {
      azureFile: {
        accountName: storageAccount.name
        accessMode: 'ReadWrite'
        accountKey: storageAccount.listKeys().keys[0].value
        shareName: fileShareName
      }
    }
  }
}

resource hbbs 'Microsoft.App/containerApps@2023-05-01' = {
  name: 'hbbs'
  location: location
  properties: {
    managedEnvironmentId: containerAppEnv.id
    configuration: {
      ingress: {
        external: true
        targetPort: 21116
        transport: 'tcp'
      }
      secrets: []
    }
    template: {
      scale: {
        minReplicas: 1
        maxReplicas: 1
      }
      containers: [{
        name: 'hbbs'
        image: 'rustdesk/rustdesk-server:latest'
        command: ['hbbs']
        volumeMounts: [{
          mountPath: '/root'
          volumeName: 'data-volume'
        }]
        resources: {
          cpu: any('0.25')
          memory: '0.5Gi'
        }
      }]
      volumes: [{
        name: 'data-volume'
        storageType: 'AzureFile'
        storageName: 'azurefiles'
        mountOptions: 'dir_mode=0777,file_mode=0777,uid=0,gid=0,mfsymlinks,nobrl,cache=strict'
      }]
    }
  }
}

resource hbbr 'Microsoft.App/containerApps@2023-05-01' = {
  name: 'hbbr'
  location: location
  properties: {
    managedEnvironmentId: containerAppEnv.id
    configuration: {
      ingress: {
        external: true
        targetPort: 21117
        transport: 'tcp'
      }
      secrets: []
    }
    template: {
      scale: {
        minReplicas: 1
        maxReplicas: 1
      }
      containers: [{
        name: 'hbbr'
        image: 'rustdesk/rustdesk-server:latest'
        command: ['hbbr']
        volumeMounts: [{
          mountPath: '/root'
          volumeName: 'data-volume'
        }]
        resources: {
          cpu: any('0.25')
          memory: '0.5Gi'
        }
      }]
      volumes: [{
        name: 'data-volume'
        storageType: 'AzureFile'
        storageName: 'azurefiles'
        mountOptions: 'dir_mode=0777,file_mode=0777,uid=0,gid=0,mfsymlinks,nobrl,cache=strict'
      }]
    }
  }
}


// 192.168.4.153 - UXAF6L2W4uJlVRY90SnSrmqdaio6O+Xd+PePYQbzDk4=
// 23.96.241.171 -  WsVrl83CvxG98HlFdUPLaYF23VVjJrZO0IVbp5+kIzE=
