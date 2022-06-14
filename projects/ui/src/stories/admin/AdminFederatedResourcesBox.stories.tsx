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
  createFederatedVirtualService,
  createListFederatedRouteTables,
} from 'stories/mocks/generators';
import { MemoryRouter } from 'react-router';
import {
  FederatedGateway,
  FederatedRouteTable,
  FederatedVirtualService,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import {
  FederatedSettings,
  FederatedUpstream,
  FederatedUpstreamGroup,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';
import { FederatedAuthConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/bootstrap_pb_service';

export default {
  title: 'Admin / AdminFederatedResourcesBox',
  component: AdminFederatedResourcesBox,
} as ComponentMeta<typeof AdminFederatedResourcesBox>;

type AdminFederatedResourcesBoxType = {
  error?: ServiceError;
  isValidating?: boolean;
  mutate?: any;
  federatedData?: FederatedVirtualService.AsObject[];
  listFederatedRouteTablesData?: FederatedRouteTable.AsObject[];
  listFederatedUpstreamsData?: FederatedUpstream.AsObject[];
  listFederatedUpstreamGroupsData?: FederatedUpstreamGroup.AsObject[];
  listFederatedAuthConfigsData?: FederatedAuthConfig.AsObject[];
  listFederatedGatewaysData?: FederatedGateway.AsObject[];
  listFederatedSettingsData?: FederatedSettings.AsObject[];
};

// TODO:  Update this with proper args
const Template: ComponentStory<typeof AdminFederatedResourcesBox> = (
  args: AdminFederatedResourcesBoxType | any
) => {
  let newArgs: AdminFederatedResourcesBoxType = { ...args };
  const useListFederatedVirtualServicesDi = injectable(
    Apis.useListFederatedVirtualServices,
    () => {
      return {
        data: newArgs.federatedData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
    }
  );

  const useListFederatedRouteTablesDi = injectable(
    Apis.useListFederatedRouteTables,
    () => {
      return {
        data: newArgs.listFederatedRouteTablesData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
    }
  );

  const useListFederatedUpstreamsDi = injectable(
    Apis.useListFederatedUpstreams,
    () => {
      return {
        data: newArgs.listFederatedUpstreamsData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
    }
  );

  const useListFederatedUpstreamGroupsDi = injectable(
    Apis.useListFederatedUpstreamGroups,
    () => {
      return {
        data: newArgs.listFederatedUpstreamGroupsData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
    }
  );

  const useListFederatedAuthConfigsDi = injectable(
    Apis.useListFederatedAuthConfigs,
    () => {
      return {
        data: newArgs.listFederatedAuthConfigsData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
    }
  );

  const useListFederatedGatewaysDi = injectable(
    Apis.useListFederatedGateways,
    () => {
      return {
        data: newArgs.listFederatedGatewaysData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
    }
  );

  const useListFederatedSettingsDi = injectable(
    Apis.useListFederatedSettings,
    () => {
      return {
        data: newArgs.listFederatedSettingsData,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating ?? false,
      };
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

const federatedData = Array.from({ length: 2 }).map(() => {
  return createFederatedVirtualService();
});

const listFederatedRouteTablesData = Array.from({ length: 2 }).map(() => {
  return createListFederatedRouteTables();
});

const listFederatedUpstreamsData = Array.from({ length: 2 }).map(() => {
  return createFederatedUpstream();
});

const listFederatedUpstreamGroupsData = Array.from({ length: 2 }).map(() => {
  return createFederatedUpstreamGroup();
});

const listFederatedAuthConfigsData = Array.from({ length: 2 }).map(() => {
  return createFederatedAuthConfig();
});

const listFederatedGatewaysData = Array.from({ length: 2 }).map(() => {
  return createFederatedGateway();
});

const listFederatedSettingsData = Array.from({ length: 2 }).map(() => {
  return createFederatedSettings();
});

export const Primary = Template.bind({});
Primary.args = {
  mutate: jest.fn(),
  isValidating: false,
  error: undefined,
  federatedData,
  listFederatedRouteTablesData,
  listFederatedUpstreamsData,
  listFederatedUpstreamGroupsData,
  listFederatedAuthConfigsData,
  listFederatedGatewaysData,
  listFederatedSettingsData,
} as Partial<
  typeof AdminFederatedResourcesBox & AdminFederatedResourcesBoxType
>;
Primary.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-box-summary');
  expect(summary).not.toBeNull();
};
