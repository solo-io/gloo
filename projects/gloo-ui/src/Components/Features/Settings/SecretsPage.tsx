import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { isEqual } from 'lodash';
import * as React from 'react';
import { TableActionCircle, TableActions } from 'Styles';
import { getIcon, RadioFilters } from 'utils/helpers';
import { SecretForm, SecretValuesType } from './SecretForm';
import { Popconfirm } from 'antd';
import { ReactComponent as KeyRing } from 'assets/key-on-ring.svg';
import useSWR from 'swr';
import { secretAPI } from 'store/secrets/api';
import { OauthSecret } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb';
import {
  AwsSecret,
  AzureSecret,
  Secret,
  TlsSecret
} from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/secret_pb';
import { useDispatch } from 'react-redux';
import { createSecret, deleteSecret } from 'store/secrets/actions';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { SoloCheckbox } from 'Components/Common/SoloCheckbox';
import styled from '@emotion/styled';
import { StyledHeader } from 'Components/Common/ListingFilter';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
interface Props {
  toggleSuccessModal?: React.Dispatch<React.SetStateAction<boolean>>;
}

const SecretsContainer = styled.div`
  display: grid;
  grid-template-areas:
    'header header'
    'sidebar content';
  grid-template-columns: 200px 1fr;
  grid-template-rows: auto 1fr;
  grid-column-gap: 20px;
`;

const FilterHeader = styled.div`
  ${StyledHeader};
  width: 185px;
`;

export const SecretsPage = (props: Props) => {
  const { data: secretsList, error } = useSWR(
    'listSecrets',
    secretAPI.getSecretsList
  );
  const dispatch = useDispatch();

  const [awsSecrets, setAwsSecrets] = React.useState<Secret.AsObject[]>([]);
  const [azureSecrets, setAzureSecrets] = React.useState<Secret.AsObject[]>([]);
  const [tlsSecrets, setTlsSecrets] = React.useState<Secret.AsObject[]>([]);
  const [oAuthSecrets, setOAuthSecrets] = React.useState<Secret.AsObject[]>([]);
  React.useEffect(() => {
    if (secretsList && secretsList.length) {
      setAwsSecrets(secretsList.filter(s => !!s.aws));
      setAzureSecrets(secretsList.filter(s => !!s.azure));
      setOAuthSecrets(secretsList.filter(s => !!s.oauth));
      setTlsSecrets(secretsList.filter(s => !!s.tls));
    }
  }, [secretsList?.length]);

  if (!secretsList) {
    return <div>Loading...</div>;
  }
  async function handleCreateSecret(
    values: SecretValuesType,
    secretKind: Secret.KindCase
  ) {
    let newSecret = new Secret().toObject();

    const {
      secretResourceRef: { name, namespace }
    } = values;

    let aws: AwsSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.AWS) {
      aws = values.awsSecret;
      newSecret = { ...newSecret, aws };
    }

    let azure: AzureSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.AZURE) {
      azure = values.azureSecret;
      newSecret = { ...newSecret, azure };
    }

    let tls: TlsSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.TLS) {
      tls = values.tlsSecret;
      newSecret = { ...newSecret, tls };
    }

    let oauth: OauthSecret.AsObject | undefined = undefined;
    if (secretKind === Secret.KindCase.OAUTH) {
      oauth = values.oAuthSecret;
      newSecret = { ...newSecret, oauth };
    }

    dispatch(
      createSecret({
        secret: {
          ...newSecret,
          metadata: {
            ...newSecret.metadata!,
            name,
            namespace
          }
        }
      })
    );
  }

  async function handleDeleteSecret(
    name: string,
    namespace: string,
    secretKind: Secret.KindCase
  ) {
    try {
      dispatch(deleteSecret({ ref: { name, namespace } }));
    } catch (error) {
      console.error(error);
    }
  }
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
                handleDeleteSecret(
                  record.name,
                  record.namespace,
                  Secret.KindCase.TLS
                )
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
                handleDeleteSecret(
                  record.name,
                  record.namespace,
                  Secret.KindCase.OAUTH
                )
              }>
              x
            </TableActionCircle>
          </div>
        </TableActions>
      )
    }
  ];

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
                handleDeleteSecret(
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
                handleDeleteSecret(
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
      <div>
        <div>
          <SectionCard cardName={'AWS Secrets'} logoIcon={getIcon('Aws')}>
            {awsTableData.length ? (
              <SoloTable
                columns={AWSColumns}
                dataSource={awsTableData}
                formComponent={() => (
                  <SecretForm
                    secretKind={Secret.KindCase.AWS}
                    onCreateSecret={handleCreateSecret}
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
                  onCreateSecret={handleCreateSecret}
                />
              )}
            />
          </SectionCard>
          <SectionCard
            cardName={'TLS'}
            logoIcon={
              <span className='text-blue-500'>
                <KeyRing className='fill-current ' />
              </span>
            }>
            {tlsTableData.length ? (
              <SoloTable
                columns={TLSColumns}
                dataSource={tlsTableData}
                formComponent={() => (
                  <SecretForm
                    secretKind={Secret.KindCase.TLS}
                    onCreateSecret={handleCreateSecret}
                  />
                )}
              />
            ) : (
              <div>No Secrets</div>
            )}
          </SectionCard>
          <SectionCard
            cardName={'OAuth'}
            logoIcon={
              <span className='text-blue-500'>
                <KeyRing className='fill-current ' />
              </span>
            }>
            {oAuthTableData.length ? (
              <SoloTable
                columns={OAuthColumns}
                dataSource={oAuthTableData}
                formComponent={() => (
                  <SecretForm
                    secretKind={Secret.KindCase.OAUTH}
                    onCreateSecret={handleCreateSecret}
                  />
                )}
              />
            ) : (
              <div>No Secrets</div>
            )}
          </SectionCard>
        </div>
      </div>
    </>
  );
};
