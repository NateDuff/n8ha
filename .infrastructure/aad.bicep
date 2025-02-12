extension microsoftGraphV1

resource clientApp 'Microsoft.Graph/applications@v1.0'  = {
  displayName: 'DCS Azure Service Manager'
  uniqueName: 'dcs-azure-admin'
  signInAudience: 'AzureADMultipleOrgs'
}

resource rootFederatedIdentity 'Microsoft.Graph/applications/federatedIdentityCredentials@v1.0' = {
  audiences: [
    'api://AzureADTokenExchange'
  ]
  description: 'Duff Cloud Services Azure DevOps Instance'
  issuer: 'https://vstoken.dev.azure.com/c121ba31-42b1-4e71-a9ec-af65377f98e1'
  name: '${clientApp.uniqueName}/DCS_ADO'
  subject: 'sc://duffcloudservices/DCS/DCS Azure Service Manager'
}

resource internalFederatedIdentity 'Microsoft.Graph/applications/federatedIdentityCredentials@v1.0' = {
  audiences: [
    'api://AzureADTokenExchange'
  ]
  description: 'Duff Cloud Services Azure DevOps Instance'
  issuer: 'https://vstoken.dev.azure.com/c121ba31-42b1-4e71-a9ec-af65377f98e1'
  name: '${clientApp.uniqueName}/DCS_ADO_Internal'
  subject: 'sc://duffcloudservices/DCS/DCS - Internal'
}

resource partnerFederatedIdentity 'Microsoft.Graph/applications/federatedIdentityCredentials@v1.0' = {
  audiences: [
    'api://AzureADTokenExchange'
  ]
  description: 'Duff Cloud Services Azure DevOps Instance'
  issuer: 'https://vstoken.dev.azure.com/c121ba31-42b1-4e71-a9ec-af65377f98e1'
  name: '${clientApp.uniqueName}/DCS_ADO_Partner'
  subject: 'sc://duffcloudservices/DCS/DCS - Partner'
}
