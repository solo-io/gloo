import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminInnerPagesWrapper } from '../../Components/Features/Admin/AdminInnerPagesWrapper';
import { MemoryRouter, useNavigate, useParams } from 'react-router';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { expect, jest } from '@storybook/jest';
import { within } from '@storybook/testing-library';
import {
  createClusterObjectRef,
  createFederatedVirtualService,
} from 'stories/mocks/generators';

export default {
  title: 'Admin / AdminInnerPagesWrapper',
  component: AdminInnerPagesWrapper,
} as ComponentMeta<typeof AdminInnerPagesWrapper>;

enum ADMIN_PAGE {
  VIRTUAL_SERVICES = 'virtual-services',
  UPSTREAMS = 'upstreams',
  UPSTREAM_GROUPS = 'upstream-groups',
  AUTHORIZATIONS = 'authorizations',
  RATE_LIMITS = 'rate-limits',
  ROUTE_TABLES = 'route-tables',
  GATEWAYS = 'gateways',
  SETTINGS = 'settings',
}

const Template: ComponentStory<typeof AdminInnerPagesWrapper> = (args: any) => {
  const useParamsDi = injectable(useParams, () => {
    return {
      adminPage: ADMIN_PAGE.VIRTUAL_SERVICES,
    };
  });

  const useListFederatedVirtualServicesDi = injectable(
    Apis.useListFederatedVirtualServices,
    () => {
      return {
        data: args.federatedData,
        error: undefined,
        isValidating: false,
        mutate: jest.fn(),
      };
    }
  );
  return (
    <DiProvider use={[useParamsDi, useListFederatedVirtualServicesDi]}>
      <MemoryRouter>
        <AdminInnerPagesWrapper />
      </MemoryRouter>
    </DiProvider>
  );
};

export const VirtualServicesWrapper = Template.bind({});
const federatedData = Array.from({ length: 2 }).map(() => {
  return createFederatedVirtualService();
});
VirtualServicesWrapper.args = {
  cardName: AdminInnerPagesWrapper.name,
  federatedData,
} as Partial<typeof AdminInnerPagesWrapper>;
VirtualServicesWrapper.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = canvas.getByTestId('admin-inner-pages-wrapper');
  expect(summary).not.toBeNull();
};
