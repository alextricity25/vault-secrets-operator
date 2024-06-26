// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package credentials

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/hashicorp/vault-secrets-operator/api/v1beta1"
	"github.com/hashicorp/vault-secrets-operator/internal/credentials/hcp"
	"github.com/hashicorp/vault-secrets-operator/internal/credentials/provider"
	"github.com/hashicorp/vault-secrets-operator/internal/credentials/vault"
	"github.com/hashicorp/vault-secrets-operator/internal/credentials/vault/consts"
)

var ProviderMethodsSupported = []string{
	consts.ProviderMethodKubernetes,
	consts.ProviderMethodJWT,
	consts.ProviderMethodAppRole,
	consts.ProviderMethodAWS,
	consts.ProviderMethodGCP,
	hcp.ProviderMethodServicePrincipal,
}

func NewCredentialProvider(ctx context.Context, client client.Client, obj client.Object, providerNamespace string) (provider.CredentialProviderBase, error) {
	var p provider.CredentialProviderBase
	switch authObj := obj.(type) {
	case *v1beta1.VaultAuth:
		var prov vault.CredentialProvider
		switch authObj.Spec.Method {
		case consts.ProviderMethodJWT:
			prov = &vault.JWTCredentialProvider{}
		case consts.ProviderMethodAppRole:
			prov = &vault.AppRoleCredentialProvider{}
		case consts.ProviderMethodKubernetes:
			prov = &vault.KubernetesCredentialProvider{}
		case consts.ProviderMethodAWS:
			prov = &vault.AWSCredentialProvider{}
		case consts.ProviderMethodGCP:
			prov = &vault.GCPCredentialProvider{}
		default:
			return nil, fmt.Errorf("unsupported authentication method %s", authObj.Spec.Method)
		}

		if err := prov.Init(ctx, client, authObj, providerNamespace); err != nil {
			return nil, err
		}

		p = prov
	case *v1beta1.HCPAuth:
		var prov hcp.CredentialProviderHCP
		switch authObj.Spec.Method {
		case hcp.ProviderMethodServicePrincipal:
			prov = &hcp.ServicePrincipleCredentialProvider{}
		default:
			return nil, fmt.Errorf("unsupported authentication method %s", authObj.Spec.Method)
		}

		if err := prov.Init(ctx, client, authObj, providerNamespace); err != nil {
			return nil, err
		}

		p = prov
	default:
		return nil, fmt.Errorf("unsupported auth object %T", authObj)
	}
	return p, nil
}
