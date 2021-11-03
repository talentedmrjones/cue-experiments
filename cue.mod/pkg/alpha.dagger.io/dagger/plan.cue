package dagger

import (
  "struct"
  sec "alpha.dagger.io/dagger/secrets"
)

_#arch: "arm" | "arm64" | "amd64"

#Plan: {
  environment: string
  description?: string
  arch: _#arch | [..._#arch]

  // Access to core system capabilities
  system: {
    secrets: [secret=string]: sec.#Provider
    // FS: API to securely access host filesystem
    fs: [dirname=string]: string
    // Proxy: API to publish and consume network services
    proxy: [servicename=string]: {...}
    // env vars
    env: [name=string]: string // eg: awsAccessKey: "AWS_ACCESS_KEY"
  }


  runtime: {
    // flesh out with all available fields:types
  }

  // Containerized actions
  actions: [name=string]: {...} & struct.MinFields(1)
}