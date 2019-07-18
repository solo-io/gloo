import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { ReactComponent as KeyRing } from 'assets/key-on-ring.svg';
import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { SecretForm } from './SecretForm';

const TLSColumns = [
  {
    title: 'Name',
    dataIndex: 'name',
    width: 200
  },
  {
    title: 'Namespace',
    dataIndex: 'namespace'
  },
  {
    title: 'Cert Chain',
    dataIndex: 'certChain'
  },
  {
    title: 'TLS Private Key',
    dataIndex: 'tlsPrivateKey'
  },
  {
    title: 'Root CA',
    dataIndex: 'rootCA'
  },
  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text: any) => <div>ACTION!</div>
  }
];

const OAuthColumns = [
  {
    title: 'Name',
    dataIndex: 'name',
    width: 200
  },
  {
    title: 'Namespace',
    dataIndex: 'namespace'
  },
  {
    title: 'Client Secret',
    dataIndex: 'clientSecret'
  },
  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text: any) => <div>ACTION!</div>
  }
];

interface Props {
  tlsSecrets?: Secret.AsObject[];
  oAuthSecrets?: Secret.AsObject[];
}

export const SecurityPage: React.FunctionComponent<Props> = props => {
  const { tlsSecrets, oAuthSecrets } = props;

  let tlsTableData: any[] = [];
  if (tlsSecrets) {
    tlsTableData = tlsSecrets.map(tlsSecret => {
      return {
        key: `${tlsSecret.metadata!.name}-${tlsSecret.metadata!.namespace}`,
        name: tlsSecret.metadata!.name,
        namespace: tlsSecret.metadata!.namespace,
        certChain: '(cert chain)',
        tlsPrivateKey: '**************************',
        rootCA: '(root ca)',
        actions: ``
      };
    });
    // This is to be replaced by the add new row form
    tlsTableData.push({
      key: '',
      name: '',
      namespace: '',
      certChain: '(cert chain)',
      tlsPrivateKey: '**************************',
      rootCA: '(root ca)',
      actions: ``
    });
  }

  let oAuthTableData: any[] = [];
  if (oAuthSecrets) {
    oAuthTableData = oAuthSecrets.map(oAuthSecret => {
      return {
        key: `${oAuthSecret.metadata!.name}-${oAuthSecret.metadata!.namespace}`,
        name: oAuthSecret.metadata!.name,
        namespace: oAuthSecret.metadata!.namespace,
        clientSecret: '(client secret)',
        actions: ``
      };
    });
    // This is to be replaced by the add new row form
    oAuthTableData.push({
      key: '',
      name: '',
      namespace: '',
      clientSecret: '(client secret)',
      actions: ``
    });
  }

  return (
    <React.Fragment>
      <SectionCard cardName={'TLS'} logoIcon={<KeyRing />}>
        {tlsTableData.length ? (
          <SoloTable
            columns={TLSColumns}
            dataSource={tlsTableData}
            formComponent={() => (
              <SecretForm secretKind={Secret.KindCase.TLS} />
            )}
          />
        ) : (
          <div>No Secrets</div>
        )}
      </SectionCard>
      <SectionCard cardName={'OAuth'} logoIcon={<KeyRing />}>
        {oAuthTableData.length ? (
          <SoloTable
            columns={OAuthColumns}
            dataSource={oAuthTableData}
            formComponent={() => (
              <SecretForm secretKind={Secret.KindCase.EXTENSION} />
            )}
          />
        ) : (
          <div>No Secrets</div>
        )}
      </SectionCard>
    </React.Fragment>
  );
};
