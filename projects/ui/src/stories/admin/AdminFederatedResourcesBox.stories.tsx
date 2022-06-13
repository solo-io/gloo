import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminFederatedResourcesBox } from '../../Components/Features/Admin/AdminBoxSummary';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { expect, jest } from '@storybook/jest';
import * as Apis from 'API/hooks';
import { within } from '@storybook/testing-library';
import {
  createFederatedAuthConfig,
  createFederatedGateway,
  createFederatedSettings,
  createFederatedUpstream,
  createFederatedUpstreamGroup,
  createFederatedUpstreamGroupSpec,
  createFederatedVirtualService,
  createFederatedVirtualServiceSpec,
  createListFederatedRouteTables,
} from 'stories/mocks/generators';
import { MemoryRouter } from 'react-router';

export default {
  title: `Admin / ${AdminFederatedResourcesBox.name}`,
  component: AdminFederatedResourcesBox,
} as unknown as ComponentMeta<typeof AdminFederatedResourcesBox>;

const Template: ComponentStory<typeof AdminFederatedResourcesBox> = args => {
  const useListFederatedVirtualServicesDi = injectable(
    Apis.useListFederatedVirtualServices,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createFederatedVirtualService();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  const useListFederatedRouteTablesDi = injectable(
    Apis.useListFederatedRouteTables,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createListFederatedRouteTables();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  const useListFederatedUpstreamsDi = injectable(
    Apis.useListFederatedUpstreams,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createFederatedUpstream();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  const useListFederatedUpstreamGroupsDi = injectable(
    Apis.useListFederatedUpstreamGroups,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createFederatedUpstreamGroup();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  const useListFederatedAuthConfigsDi = injectable(
    Apis.useListFederatedAuthConfigs,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createFederatedAuthConfig();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  const useListFederatedGatewaysDi = injectable(
    Apis.useListFederatedGateways,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createFederatedGateway();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  const useListFederatedSettingsDi = injectable(
    Apis.useListFederatedSettings,
    () => {
      const data = Array.from({ length: 2 }).map(() => {
        return createFederatedSettings();
      });
      return { data, error: undefined, mutate: jest.fn(), isValidating: false };
    }
  );

  return (
    <DiProvider
      use={[
        useListFederatedVirtualServicesDi,
        useListFederatedRouteTablesDi,
        useListFederatedUpstreamsDi,
        useListFederatedUpstreamGroupsDi,
        useListFederatedAuthConfigsDi,
        useListFederatedGatewaysDi,
        useListFederatedSettingsDi,
      ]}>
      <MemoryRouter>
        <AdminFederatedResourcesBox />
      </MemoryRouter>
    </DiProvider>
  );
};

export const Primary = Template.bind({});
Primary.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-box-summary');
  expect(summary).not.toBeNull();
};
Primary.args = {} as Partial<typeof AdminFederatedResourcesBox>;
