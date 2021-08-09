import React from 'react';
import styled from '@emotion/styled';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import {
  getFunctionList,
  getUpstreamType,
  TYPE_AWS,
  TYPE_AZURE,
  TYPE_KUBE,
  TYPE_STATIC,
} from 'utils/upstream-helpers';
import {
  CardSubsectionWrapper,
  CardSubsectionContent,
} from 'Components/Common/Card';
import { Loading } from 'Components/Common/Loading';
import { SoloInput } from 'Components/Common/SoloInput';
import { StringCardsList } from 'Components/Common/StringCardsList';

const Wrapper = styled(CardSubsectionWrapper)`
  display: grid;
  grid-template-columns: 1fr 2fr;
`;

const SectionTitle = styled.div`
  margin-bottom: 8px;
`;

type ConfigProps = {
  cols: number;
};
const ConfigContainer = styled.div<ConfigProps>`
  flex: ${(props: ConfigProps) => props.cols};
`;
const ConfigSection = styled(CardSubsectionContent)<ConfigProps>`
  display: grid;
  grid-template-columns: repeat(${(props: ConfigProps) => props.cols}, 1fr);
  grid-gap: 18px 10px;
  height: auto;
`;
const SecurityContainer = styled.div`
  height: 100%;
  display: grid;
  grid-template-rows: 27px 1fr;
  margin-left: 18px;
`;
const SecuritySection = styled(CardSubsectionContent)`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  grid-gap: 18px 10px;
  height: auto;
`;
const FunctionsContainer = styled.div`
  position: relative;
  grid-column: span 3;
  margin-top: 18px;
`;
const FunctionsSection = styled(CardSubsectionContent)`
  height: auto;
`;

const ListBlock = styled.div`
  margin: 8px 0 5px;
  .ant-list {
    background: white;
    
    .ant-list-item {
      padding: 0 5px;
    }
  }
`;

type Field = {
  title: string;
  value?: boolean | number | string;
};

type Config = {
  fields: Field[];
  cols?: number;
};

// return the config values to display, depending on the upstream type
const getUpstreamTypeConfig = (
  upstream?: Upstream.AsObject
): Config | undefined => {
  const upstreamType = getUpstreamType(upstream);
  switch (upstreamType) {
    case TYPE_AWS: {
      return {
        fields: [
          {
            title: 'Region',
            value: upstream?.spec?.aws?.region,
          },
          {
            title: 'Secret Name',
            value: upstream?.spec?.aws?.secretRef?.name,
          },
          {
            title: 'Secret Namespace',
            value: upstream?.spec?.aws?.secretRef?.namespace,
          },
        ],
      };
    }
    case TYPE_AZURE: {
      return {
        fields: [
          {
            title: 'Function App Name',
            value: upstream?.spec?.azure?.functionAppName,
          },
          {
            title: 'Secret Name',
            value: upstream?.spec?.azure?.secretRef?.name,
          },
          {
            title: 'Secret Namespace',
            value: upstream?.spec?.azure?.secretRef?.namespace,
          },
        ],
      };
    }
    case TYPE_KUBE: {
      return {
        fields: [
          {
            title: 'Service Name',
            value: upstream?.spec?.kube?.serviceName,
          },
          {
            title: 'Service Namespace',
            value: upstream?.spec?.kube?.serviceNamespace,
          },
          {
            title: 'Service Port',
            value: upstream?.spec?.kube?.servicePort,
          },
        ],
      };
    }
    case TYPE_STATIC: {
      const hostsList = upstream?.spec?.pb_static?.hostsList;
      const fields: Field[] = hostsList
        ? hostsList
            .map(host => [
              { title: 'Host Address', value: host.addr },
              { title: 'Port', value: host.port },
            ])
            .flat()
        : [];
      fields.push({
        title: 'Use TLS',
        value: upstream?.spec?.pb_static?.useTls,
      });
      return { fields, cols: 2 };
    }
    default:
      return undefined;
  }
};

type Props = {
  upstream?: Upstream.AsObject;
};

const UpstreamConfiguration = ({ upstream }: Props) => {
  const upstreamConfig = getUpstreamTypeConfig(upstream);
  const functionsList = getFunctionList(upstream);
  const cols = upstreamConfig?.cols || 1;

  return upstream ? (
    <Wrapper>
      {upstreamConfig && (
        <ConfigContainer cols={cols}>
          <SectionTitle>{getUpstreamType(upstream)} Settings</SectionTitle>
          <ConfigSection cols={cols}>
            {upstreamConfig.fields.map(({ title, value }, idx) => (
              <SoloInput
                key={`${title}-${idx}`}
                title={title}
                value={(value ?? '').toString()}
                disabled
              />
            ))}
          </ConfigSection>
        </ConfigContainer>
      )}
      <SecurityContainer>
        <SectionTitle>Security</SectionTitle>
        <SecuritySection>
          <SoloInput
            title='TLS Certificate'
            value={upstream?.spec?.sslConfig?.sslFiles?.tlsCert ?? ''}
            disabled
          />
          <SoloInput
            title='TLS Private Key'
            value={upstream?.spec?.sslConfig?.sslFiles?.tlsKey ?? ''}
            disabled
          />
          <SoloInput
            title='Root Certificate'
            value={upstream?.spec?.sslConfig?.sslFiles?.rootCa ?? ''}
            disabled
          />
        </SecuritySection>
      </SecurityContainer>
      {functionsList.length > 0 && (
        <FunctionsContainer>
          <SectionTitle>Functions</SectionTitle>
          <FunctionsSection>
            <StringCardsList values={functionsList} />
          </FunctionsSection>
        </FunctionsContainer>
      )}
    </Wrapper>
  ) : (
    <Loading message='Retrieving upstream details...' />
  );
};

export default UpstreamConfiguration;
