import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { ReactComponent as KeyRing } from 'assets/key-on-ring.svg';

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
    title: 'Key/Value',
    dataIndex: 'keyValue'
  },
  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text: any) => <div>ACTION!</div>
  }
];

interface Props {}

export const SecretsPage = (props: Props) => {
  const awsData: any[] = [];
  for (let i = 1; i <= 10; i++) {
    awsData.push({
      key: i,
      name: 'John Brown',
      namespace: `${i}2`,
      accessKey: `New York No. ${i} Lake Park`,
      secretKey: `My name is John Brown, I am ${i}2 years old, living in New York No. ${i} Lake Park.`,
      actions: ``
    });
  }

  const oAuthData: any[] = [];
  for (let i = 1; i <= 10; i++) {
    oAuthData.push({
      key: i,
      name: 'John Brown',
      namespace: `${i}2`,
      keyValue: `New York No. ${i} Lake Park`,
      actions: ``
    });
  }

  return (
    <React.Fragment>
      <SectionCard cardName={'AWS Secrets'} logoIcon={<KeyRing />}>
        <SoloTable columns={AWSColumns} dataSource={awsData} />
      </SectionCard>
      <SectionCard cardName={'OAuth'} logoIcon={<KeyRing />}>
        <SoloTable columns={OAuthColumns} dataSource={oAuthData} />
      </SectionCard>
    </React.Fragment>
  );
};
