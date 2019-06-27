import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloTable } from 'Components/Common/SoloTable';
import { ReactComponent as KeyRing } from 'assets/key-on-ring.svg';

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

interface Props {}

export const SecurityPage = (props: Props) => {
  const tlsData: any[] = [];
  for (let i = 1; i <= 10; i++) {
    tlsData.push({
      key: i,
      name: 'John Brown',
      namespace: `${i}2`,
      certChain: `New York No. ${i} Lake Park`,
      tlsPrivateKey: `My name is John Brown, I am ${i}2 years old, living in New York No. ${i} Lake Park.`,
      rootCA: `New York No. ${i} Lake Park`,
      actions: ``
    });
  }

  const oAuthData: any[] = [];
  for (let i = 1; i <= 10; i++) {
    oAuthData.push({
      key: i,
      name: 'John Brown',
      namespace: `${i}2`,
      clientSecret: `New York No. ${i} Lake Park`,
      actions: ``
    });
  }

  return (
    <React.Fragment>
      <SectionCard cardName={'TLS'} logoIcon={<KeyRing />}>
        <SoloTable columns={TLSColumns} dataSource={tlsData} />
      </SectionCard>
      <SectionCard cardName={'OAuth'} logoIcon={<KeyRing />}>
        <SoloTable columns={OAuthColumns} dataSource={oAuthData} />
      </SectionCard>
    </React.Fragment>
  );
};
