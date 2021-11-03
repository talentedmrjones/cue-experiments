package secrets

#Provider: {sops: _#sops} | {gpg: _#gpg} | {"aws-vault": _#awsVault}

_#sops: {
  file: string    // path to encrypted sops file
  config?: string // path to sops config file

  // when runtime fills
  [string]: string
}

_#awsVault: {
  profile: string
}

_#gpg: {
  passprompt: *false | true
}