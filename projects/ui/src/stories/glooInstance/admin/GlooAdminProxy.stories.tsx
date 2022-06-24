import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminProxy } from '../../../Components/Features/GlooInstance/Admin/GlooAdminProxy';
import { MemoryRouter, useParams } from 'react-router';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { faker } from '@faker-js/faker';
import { createProxy } from 'stories/mocks/generators';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { Proxy } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import { glooResourceApi } from 'API/gloo-resource';

export default {
  title: 'GlooInstance / Admin / GlooAdminProxy',
  component: GlooAdminProxy,
} as ComponentMeta<typeof GlooAdminProxy>;

interface GlooAdminProxyStoryProps {
  name: string;
  namespace: string;
  data: Proxy.AsObject[];
  error?: any;
  isValidating: boolean;
  mutate: any;
  yaml: string;
}

const Template: ComponentStory<typeof GlooAdminProxy> = (args: any) => {
  const newArgs: GlooAdminProxyStoryProps = args;
  const useParamsDi = injectable(useParams, () => {
    return {
      name: newArgs.name,
      namespace: newArgs.namespace,
    };
  });

  const useListProxiesDi = injectable(Apis.useListProxies, () => {
    return {
      data: newArgs.data,
      error: newArgs.error,
      isValidating: newArgs.isValidating,
      mutate: newArgs.mutate,
    };
  });

  const getProxyYAMLDi = injectable(glooResourceApi.getProxyYAML, () => {
    return Promise.resolve(newArgs.yaml);
  });
  /**
   * useParams();
  const { data: proxies, error: pError } = useListProxies({ name, namespace });
   */
  return (
    <DiProvider use={[useParamsDi, useListProxiesDi, getProxyYAMLDi]}>
      <MemoryRouter>
        <GlooAdminProxy />
      </MemoryRouter>
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createProxy();
});

const yaml = faker.random.words();

export const Primary = Template.bind({});
Primary.args = {
  data,
  error: undefined,
  isValidating: false,
  name: faker.random.word(),
  namespace: faker.random.word(),
  mutate: jest.fn(),
  yaml,
} as Partial<typeof GlooAdminProxy>;
