import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors, TableActions, TableActionCircle } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { ReactComponent as KeyRing } from 'assets/key-on-ring.svg';
import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { SecretForm } from './SecretForm';
import { DeleteSecretRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { useDeleteSecret } from 'Api';

interface Props {
  tlsSecrets?: Secret.AsObject[];
  oAuthSecrets?: Secret.AsObject[];
}

export const SecurityPage: React.FunctionComponent<Props> = props => {
  const { tlsSecrets, oAuthSecrets } = props;
  const [tlsSecretsList, setTlsSecretsList] = React.useState<Secret.AsObject[]>(
    !!tlsSecrets ? tlsSecrets : []
  );
  const [oAuthSecretsList, setOAuthSecretsList] = React.useState<
    Secret.AsObject[]
  >(!!oAuthSecrets ? oAuthSecrets : []);
  const { refetch: makeRequest } = useDeleteSecret(null);

  let tlsTableData: any[] = [];
  if (tlsSecretsList) {
    tlsTableData = tlsSecretsList.map(tlsSecret => {
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
  if (oAuthSecretsList) {
    oAuthTableData = oAuthSecretsList.map(oAuthSecret => {
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
      render: (text: any, record: any) => (
        <TableActions>
          <div style={{ marginLeft: '5px' }}>
            <TableActionCircle
              onClick={() =>
                deleteSecret(record.name, record.namespace, Secret.KindCase.TLS)
              }>
              x
            </TableActionCircle>
          </div>
        </TableActions>
      )
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
      render: (text: any, record: any) => (
        <TableActions>
          <div style={{ marginLeft: '5px' }}>
            <TableActionCircle
              onClick={() =>
                deleteSecret(
                  record.name,
                  record.namespace,
                  Secret.KindCase.EXTENSION
                )
              }>
              x
            </TableActionCircle>
          </div>
        </TableActions>
      )
    }
  ];

  function deleteSecret(
    name: string,
    namespace: string,
    secretKind: Secret.KindCase
  ) {
    if (secretKind === Secret.KindCase.TLS) {
      setTlsSecretsList(list => list.filter(s => s.metadata!.name !== name));
    }
    if (secretKind === Secret.KindCase.EXTENSION) {
      setOAuthSecretsList(list => list.filter(s => s.metadata!.name !== name));
    }
    let req = new DeleteSecretRequest();
    let ref = new ResourceRef();
    ref.setName(name);
    ref.setNamespace(namespace);
    req.setRef(ref);
    makeRequest(req);
  }

  return (
    <React.Fragment>
      <SectionCard cardName={'TLS'} logoIcon={<KeyRing />}>
        {tlsTableData.length ? (
          <SoloTable
            columns={TLSColumns}
            dataSource={tlsTableData}
            formComponent={() => (
              <SecretForm secretKind={Secret.KindCase.TLS} asTable />
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
              <SecretForm secretKind={Secret.KindCase.EXTENSION} asTable />
            )}
          />
        ) : (
          <div>No Secrets</div>
        )}
      </SectionCard>
    </React.Fragment>
  );
};
