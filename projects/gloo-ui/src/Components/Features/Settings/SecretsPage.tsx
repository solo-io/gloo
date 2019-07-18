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
import { getIcon } from 'utils/helpers';

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
    render: (text: any) => <div>ACTION!</div>
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
    render: (text: any) => <div>ACTION!</div>
  }
];

interface Props {
  awsSecrets?: Secret.AsObject[];
  azureSecrets?: Secret.AsObject[];
}

export const SecretsPage = (props: Props) => {
  const { awsSecrets, azureSecrets } = props;

  let awsTableData: any[] = [];
  if (awsSecrets) {
    awsTableData = awsSecrets.map(awsSecret => {
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
  if (azureSecrets) {
    azureTableData = azureSecrets.map(azureSecret => {
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

  return (
    <React.Fragment>
      <SectionCard cardName={'AWS Secrets'} logoIcon={getIcon('AWS')}>
        {awsTableData.length ? (
          <SoloTable
            columns={AWSColumns}
            dataSource={awsTableData}
            formComponent={() => (
              <SecretForm secretKind={Secret.KindCase.AWS} />
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
            <SecretForm secretKind={Secret.KindCase.AZURE} />
          )}
        />
      </SectionCard>
    </React.Fragment>
  );
};
