import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors, TableActions, TableActionCircle } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { ReactComponent as EditPencil } from 'assets/edit-pencil.svg';
import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { SecretForm } from './SecretForm';
import { getIcon } from 'utils/helpers';
import { DeleteSecretRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/secret_pb';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { useDeleteSecret } from 'Api';

interface Props {
  awsSecrets?: Secret.AsObject[];
  azureSecrets?: Secret.AsObject[];
}

export const SecretsPage = (props: Props) => {
  const { awsSecrets, azureSecrets } = props;
  const [awsSecretList, setAwsSecretList] = React.useState<Secret.AsObject[]>(
    !!awsSecrets ? awsSecrets : []
  );
  const [azureSecretList, setAzureSecretList] = React.useState<
    Secret.AsObject[]
  >(!!azureSecrets ? azureSecrets : []);

  const { refetch: makeRequest } = useDeleteSecret(null);
  let awsTableData: any[] = [];

  if (awsSecretList) {
    awsTableData = awsSecretList.map(awsSecret => {
      return {
        key: `${awsSecret.metadata!.name}-${awsSecret.metadata!.namespace}`,
        name: awsSecret.metadata!.name,
        namespace: awsSecret.metadata!.namespace,
        accessKey: awsSecret.aws!.accessKey,
        secretKey: awsSecret.aws!.secretKey
      };
    });
    // This is to be replaced by the add new row form
    awsTableData.push({
      key: '',
      name: '',
      namespace: '',
      accessKey: '',
      secretKey: ''
    });
  }

  let azureTableData: any[] = [];
  if (azureSecretList) {
    azureTableData = azureSecretList.map(azureSecret => {
      return {
        key: `${azureSecret.metadata!.name}-${azureSecret.metadata!.namespace}`,
        name: azureSecret.metadata!.name,
        namespace: azureSecret.metadata!.namespace,
        keyValue: azureSecret.azure!.apiKeysMap,
        actions: ``
      };
    });
    // This is to be replaced by the add new row form
    azureTableData.push({
      key: ``,
      name: '',
      namespace: '',
      keyValue: '',
      actions: ``
    });
  }

  const AWSColumns = [
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
      title: 'AccessKey',
      dataIndex: 'aws.accessKey'
    },
    {
      title: 'Secret Key',
      dataIndex: 'aws.secretKey'
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (text: any, record: any) => (
        <TableActions>
          <div style={{ marginLeft: '5px' }}>
            <TableActionCircle
              onClick={() =>
                deleteSecret(record.name, record.namespace, Secret.KindCase.AWS)
              }>
              x
            </TableActionCircle>
          </div>
        </TableActions>
      )
    }
  ];

  const AzureColumns = [
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
      title: 'Key/Value',
      dataIndex: 'keyValue'
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
                  Secret.KindCase.AZURE
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
    if (secretKind === Secret.KindCase.AWS) {
      setAwsSecretList(list => list.filter(s => s.metadata!.name !== name));
    }
    if (secretKind === Secret.KindCase.AZURE) {
      setAzureSecretList(list => list.filter(s => s.metadata!.name !== name));
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
      <SectionCard cardName={'AWS Secrets'} logoIcon={getIcon('AWS')}>
        {awsTableData.length ? (
          <SoloTable
            columns={AWSColumns}
            dataSource={awsTableData}
            formComponent={() => (
              <SecretForm secretKind={Secret.KindCase.AWS} asTable />
            )}
          />
        ) : (
          <div>No Secrets</div>
        )}
      </SectionCard>
      <SectionCard cardName={'Azure Secrets'} logoIcon={getIcon('Azure')}>
        <SoloTable
          columns={AzureColumns}
          dataSource={azureTableData}
          formComponent={() => (
            <SecretForm secretKind={Secret.KindCase.AZURE} asTable />
          )}
        />
      </SectionCard>
    </React.Fragment>
  );
};
