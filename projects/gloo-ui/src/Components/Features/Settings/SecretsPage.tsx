import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { isEqual } from 'lodash';
import { Secret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import * as React from 'react';
import { TableActionCircle, TableActions } from 'Styles';
import { getIcon } from 'utils/helpers';
import { SecretForm, SecretValuesType } from './SecretForm';
import { Popconfirm } from 'antd';

interface Props {
  awsSecrets?: Secret.AsObject[];
  azureSecrets?: Secret.AsObject[];
  toggleSuccessModal?: React.Dispatch<React.SetStateAction<boolean>>;
  onCreateSecret: (
    values: SecretValuesType,
    secretKind: Secret.KindCase
  ) => void;
  onDeleteSecret: (
    name: string,
    namespace: string,
    secretKind: Secret.KindCase
  ) => void;
}

function equivalentProps(
  oldProps: Readonly<Props>,
  nextProps: Readonly<Props>
): boolean {
  if (
    !!oldProps.awsSecrets !== !!nextProps.awsSecrets ||
    !!oldProps.azureSecrets !== !!nextProps.azureSecrets
  ) {
    return false;
  }
  if (!oldProps.awsSecrets || !oldProps.azureSecrets) {
    return true;
  }

  return (
    isEqual(oldProps.awsSecrets, nextProps.awsSecrets) &&
    isEqual(oldProps.azureSecrets, nextProps.azureSecrets)
  );
}

export const SecretsPage = React.memo((props: Props) => {
  const { awsSecrets, azureSecrets } = props;

  let awsTableData: any[] = [];

  if (awsSecrets) {
    awsTableData = awsSecrets.map(awsSecret => {
      return {
        key: `${awsSecret.metadata!.name}-${awsSecret.metadata!.namespace}`,
        name: awsSecret.metadata!.name,
        namespace: awsSecret.metadata!.namespace,
        accessKey: '**************************',
        secretKey: '**************************'
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
      dataIndex: 'accessKey'
    },
    {
      title: 'Secret Key',
      dataIndex: 'secretKey'
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (text: any, record: any) => (
        <TableActions>
          <div style={{ marginLeft: '5px' }}>
            <Popconfirm
              data-testid={`delete-aws-secret-${record.name}`}
              onConfirm={() =>
                props.onDeleteSecret(
                  record.name,
                  record.namespace,
                  Secret.KindCase.AWS
                )
              }
              title={'Are you sure you want to delete this secret? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
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
            <Popconfirm
              data-testid={`delete-azure-secret-${record.name}`}
              onConfirm={() =>
                props.onDeleteSecret(
                  record.name,
                  record.namespace,
                  Secret.KindCase.AZURE
                )
              }
              title={'Are you sure you want to delete this secret? '}
              okText='Yes'
              cancelText='No'>
              <TableActionCircle>x</TableActionCircle>
            </Popconfirm>
          </div>
        </TableActions>
      )
    }
  ];

  return (
    <>
      <SectionCard cardName={'AWS Secrets'} logoIcon={getIcon('Aws')}>
        {awsTableData.length ? (
          <SoloTable
            columns={AWSColumns}
            dataSource={awsTableData}
            formComponent={() => (
              <SecretForm
                secretKind={Secret.KindCase.AWS}
                onCreateSecret={props.onCreateSecret}
              />
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
            <SecretForm
              secretKind={Secret.KindCase.AZURE}
              onCreateSecret={props.onCreateSecret}
            />
          )}
        />
      </SectionCard>
    </>
  );
}, equivalentProps);
