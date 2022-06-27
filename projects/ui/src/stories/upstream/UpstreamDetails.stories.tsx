import { ComponentMeta, ComponentStory } from '@storybook/react';
import { getUpstreamDetails } from 'API/gloo-resource';
import * as Apis from 'API/hooks';
import { useIsGlooFedEnabled } from 'API/hooks';
import { UpstreamDetails } from 'Components/Features/Upstream/UpstreamDetails';
import React from 'react';
import { DiProvider, injectable } from 'react-magnetic-di';
import { MemoryRouter, useParams } from 'react-router';
import { createUpstream } from 'stories/mocks/generators';
import { jest } from '@storybook/jest';
import { faker } from '@faker-js/faker';

export default {
  title: 'Upstream / UpstreamDetails',
  component: UpstreamDetails,
} as ComponentMeta<typeof UpstreamDetails>;

const mockUpstreamName = 'fake-gloo';
const mockUpstreamNamespace = 'fake-gloo-namespace';
const useParamsDi = injectable(useParams, () => {
  const upstream = createUpstream();
  return {
    namespace: mockUpstreamNamespace,
    name: mockUpstreamName,
    upstreamClusterName: upstream.metadata?.clusterName,
    upstreamNamespace: upstream.metadata?.namespace,
    upstreamName: upstream.metadata?.name,
  } as ReturnType<typeof useParams>;
});

const useIsGlooFedEnabledDi = injectable(useIsGlooFedEnabled, () => {
  return {
    data: { enabled: true },
  } as ReturnType<typeof useIsGlooFedEnabled>;
});

const getMockUpstream = (args: any) => {
  const upstream = createUpstream({
    spec: {
      healthChecksList: [],
      aws: {
        lambdaFunctionsList: [
          {
            logicalName: 'test-lambda',
            lambdaFunctionName: 'test-lambda',
            qualifier: '1',
          },
        ],
        region: 'us-east1',
        roleArn: 'test-role',
        secretRef: {
          name: 'test-secret',
          namespace: 'test-secret-namespace',
        },
      },
      sslConfig: {
        alpnProtocolsList: [],
        sni: '',
        verifySubjectAltNameList: [],
        sslFiles: {
          rootCa: 'test-root-ca',
          tlsCert: 'test-tls-cert',
          tlsKey: 'test-tls-key',
        },
      },
    },
  });
  upstream.metadata!.name = mockUpstreamName;
  upstream.metadata!.namespace = mockUpstreamNamespace;
  if (args.clusterName) upstream.metadata!.clusterName = args.clusterName;
  return upstream;
};

const useGetUpstreamDetailsDi = injectable(
  Apis.useGetUpstreamDetails,
  args => ({
    data: { upstream: getMockUpstream(args) },
    isValidating: false,
    mutate: jest.fn(),
  })
);

const getUpstreamDetailsDi = injectable(
  getUpstreamDetails,
  args =>
    new Promise((resolve, _) => resolve({ upstream: getMockUpstream(args) }))
);

const useGetFailoverSchemeDi = injectable(Apis.useGetFailoverScheme, args => {
  const { name, namespace, clusterName } = args;
  return {
    data: {
      metadata: { name, namespace, clusterName },
      spec: {
        failoverGroupsList: [
          {
            priorityGroupList: [
              {
                cluster: faker.random.word(),
                upstreamsList: Array.from({ length: 20 }).map((_, i) => ({
                  name: 'upstream-' + i,
                  namespace: mockUpstreamNamespace,
                })),
              },
            ],
          },
          {
            priorityGroupList: [
              {
                cluster: faker.random.word(),
                upstreamsList: Array.from({ length: 4 }).map((_, i) => ({
                  name: 'upstream-' + i,
                  namespace: mockUpstreamNamespace,
                })),
              },
            ],
          },
          {
            priorityGroupList: [
              {
                cluster: faker.random.word(),
                upstreamsList: Array.from({ length: 2 }).map((_, i) => ({
                  name: 'upstream-' + i,
                  namespace: mockUpstreamNamespace,
                })),
                localityWeight: { value: 2 },
              },
              {
                cluster: faker.random.word(),
                upstreamsList: Array.from({ length: 3 }).map((_, i) => ({
                  name: 'upstream-' + i + 2,
                  namespace: mockUpstreamNamespace,
                })),
                localityWeight: { value: 5 },
              },
            ],
          },
        ],
      },
    },
  } as ReturnType<typeof Apis.useGetFailoverScheme>;
});

const Template: ComponentStory<typeof UpstreamDetails> = _args => {
  return (
    <DiProvider
      use={[
        useParamsDi,
        useIsGlooFedEnabledDi,
        useGetUpstreamDetailsDi,
        useGetFailoverSchemeDi,
        getUpstreamDetailsDi,
      ]}>
      <MemoryRouter>
        <UpstreamDetails />
      </MemoryRouter>
    </DiProvider>
  );
};

export const Primary = Template.bind({});
Primary.args = {};
