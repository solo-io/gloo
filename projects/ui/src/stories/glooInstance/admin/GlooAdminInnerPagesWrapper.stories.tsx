import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminInnerPagesWrapper } from '../../../Components/Features/GlooInstance/Admin/GlooAdminInnerPagesWrapper';
import { MemoryRouter, useNavigate } from 'react-router';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { faker } from '@faker-js/faker';

import { createGateway } from 'stories/mocks/generators';
import { useParams } from 'react-router';
import { gatewayResourceApi } from 'API/gateway-resources';

enum ADMIN_PAGE {
  GATEWAYS = 'gateways',
  PROXY = 'proxy',
  ENVOY = 'envoy',
  SETTINGS = 'settings',
  WATCHED_NAMESPACES = 'watched-namespaces',
  SECRETS = 'secrets',
  DOES_NOT_EXIST = 'does-not-exist',
}

export default {
  title: 'GlooInstance / Admin / GlooAdminInnerPagesWrapper',
  component: GlooAdminInnerPagesWrapper,
} as ComponentMeta<typeof GlooAdminInnerPagesWrapper>;

const Template: ComponentStory<typeof GlooAdminInnerPagesWrapper> = (
  args: any
) => {
  //   di(useParams, useNavigate, usePageGlooInstance);
  const useParamsDi = injectable(useParams, () => {
    return {
      adminPage: ADMIN_PAGE.GATEWAYS,
      name: args.name,
    };
  });
  const getGatewayYAMLDi = injectable(gatewayResourceApi.getGatewayYAML, () => {
    return Promise.resolve(faker.random.words(8));
  });

  const useListGatewaysDi = injectable(Apis.useListGateways, () => {
    return {
      data: args.data,
      error: args.error,
      mutate: args.mutate,
      isValidating: args.isValidating,
    };
  });
  return (
    <DiProvider use={[useParamsDi, useListGatewaysDi, getGatewayYAMLDi]}>
      <MemoryRouter>
        <GlooAdminInnerPagesWrapper />
      </MemoryRouter>
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createGateway();
});

export const Primary = Template.bind({});
Primary.args = {
  name: faker.random.word(),
  useNavigate: jest.fn(),
  data,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof GlooAdminInnerPagesWrapper>;
