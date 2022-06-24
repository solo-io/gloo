import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminGateways } from '../../../Components/Features/GlooInstance/Admin/GlooAdminGateways';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import * as Apis from 'API/hooks';
import { jest } from '@storybook/jest';
import { faker } from '@faker-js/faker';
import { useParams } from 'react-router';
import { glooResourceApi } from 'API/gloo-resource';
import { createGateway } from 'stories/mocks/generators';
import { gatewayResourceApi } from 'API/gateway-resources';

export default {
  title: 'GlooInstance / Admin / GlooAdminGateways',
  component: GlooAdminGateways,
} as ComponentMeta<typeof GlooAdminGateways>;

const Template: ComponentStory<typeof GlooAdminGateways> = (args: any) => {
  const useParamsDi = injectable(useParams, () => {
    return {
      name: args.name,
      namespace: args.namespace,
    };
  });
  const useListGatewaysDi = injectable(Apis.useListGateways, () => {
    return {
      data: args.data,
      error: args.error,
      mutate: args.mutate,
      isValidating: args.isValidating,
    };
  });
  const getGatewayYAMLDi = injectable(gatewayResourceApi.getGatewayYAML, () => {
    return Promise.resolve(faker.random.words(8));
  });
  return (
    <DiProvider use={[useParamsDi, useListGatewaysDi, getGatewayYAMLDi]}>
      <GlooAdminGateways />
    </DiProvider>
  );
};

const name = faker.random.word();
const namespace = faker.random.word();

const data = Array.from({ length: 1 }).map(() => {
  return createGateway();
});

export const Primary = Template.bind({});
Primary.args = {
  name,
  namespace,
  data,
  mutate: jest.fn(),
  isValidating: false,
  error: undefined,
} as Partial<typeof GlooAdminGateways>;
