import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminLanding } from '../../Components/Features/Admin/AdminLanding';
import { MemoryRouter } from 'react-router';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { expect, jest } from '@storybook/jest';
import { within } from '@storybook/testing-library';
import { createClusterDetailsObj } from 'stories/mocks/generators';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { grpc } from '@improbable-eng/grpc-web';
import { ClusterDetails } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import {
  FederatedGateway,
  FederatedRouteTable,
  FederatedVirtualService,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';
import { FederatedAuthConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import {
  FederatedUpstream,
  FederatedUpstreamGroup,
  FederatedSettings,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gloo_resources_pb';

type PartialList = {
  error?: ServiceError;
  mutate?: any;
  isValidating: boolean;
};

type AdminLandingStoryProps = {
  listDetails: PartialList & {
    data?: ClusterDetails.AsObject[];
  };
  federatedData: PartialList & {
    data: FederatedVirtualService.AsObject[];
  };
  listFederatedRouteTablesData: PartialList & {
    data: FederatedRouteTable.AsObject[];
  };
  listFederatedUpstreamsData: PartialList & {
    data: FederatedUpstream.AsObject[];
  };
  listFederatedUpstreamGroupsData: PartialList & {
    data: FederatedUpstreamGroup.AsObject[];
  };
  listFederatedAuthConfigsData: PartialList & {
    data: FederatedAuthConfig.AsObject[];
  };
  listFederatedGatewaysData: PartialList & {
    data: FederatedGateway.AsObject[];
  };
  listFederatedSettingsData: PartialList & {
    data: FederatedSettings.AsObject[];
  };
};

export default {
  title: 'Admin / AdminLanding',
  component: AdminLanding,
} as ComponentMeta<typeof AdminLanding>;

const Template: ComponentStory<typeof AdminLanding> = (args: any) => {
  const newArgs: AdminLandingStoryProps = { ...args };
  const useListClusterDetailsDi = injectable(Apis.useListClusterDetails, () => {
    return {
      data: newArgs.listDetails.data,
      error: newArgs.listDetails.error,
      mutate: newArgs.listDetails.mutate,
      isValidating: newArgs.listDetails.isValidating,
    };
  });
  const useListFederatedVirtualServicesDi = injectable(
    Apis.useListFederatedVirtualServices,
    () => {
      return {
        data: newArgs.federatedData.data,
        error: newArgs.federatedData.error,
        mutate: newArgs.federatedData.mutate,
        isValidating: newArgs.federatedData.isValidating,
      };
    }
  );

  const useListFederatedRouteTablesDi = injectable(
    Apis.useListFederatedRouteTables,
    () => {
      return {
        data: newArgs.listFederatedRouteTablesData.data,
        error: newArgs.listFederatedRouteTablesData.error,
        mutate: newArgs.listFederatedRouteTablesData.mutate,
        isValidating: newArgs.listFederatedRouteTablesData.isValidating,
      };
    }
  );

  const useListFederatedUpstreamsDi = injectable(
    Apis.useListFederatedUpstreams,
    () => {
      return {
        data: newArgs.listFederatedUpstreamsData.data,
        error: newArgs.listFederatedUpstreamsData.error,
        mutate: newArgs.listFederatedUpstreamsData.mutate,
        isValidating: newArgs.listFederatedUpstreamsData.isValidating,
      };
    }
  );

  const useListFederatedUpstreamGroupsDi = injectable(
    Apis.useListFederatedUpstreamGroups,
    () => {
      return {
        data: newArgs.listFederatedUpstreamGroupsData.data,
        error: newArgs.listFederatedUpstreamGroupsData.error,
        mutate: newArgs.listFederatedUpstreamGroupsData.mutate,
        isValidating: newArgs.listFederatedUpstreamGroupsData.isValidating,
      };
    }
  );

  const useListFederatedAuthConfigsDi = injectable(
    Apis.useListFederatedAuthConfigs,
    () => {
      return {
        data: newArgs.listFederatedAuthConfigsData.data,
        error: newArgs.listFederatedAuthConfigsData.error,
        mutate: newArgs.listFederatedAuthConfigsData.mutate,
        isValidating: newArgs.listFederatedAuthConfigsData.isValidating,
      };
    }
  );

  const useListFederatedGatewaysDi = injectable(
    Apis.useListFederatedGateways,
    () => {
      return {
        data: newArgs.listFederatedGatewaysData.data,
        error: newArgs.listFederatedGatewaysData.error,
        mutate: newArgs.listFederatedGatewaysData.mutate,
        isValidating: newArgs.listFederatedGatewaysData.isValidating,
      };
    }
  );

  const useListFederatedSettingsDi = injectable(
    Apis.useListFederatedSettings,
    () => {
      return {
        data: newArgs.listFederatedSettingsData.data,
        error: newArgs.listFederatedSettingsData.error,
        mutate: newArgs.listFederatedSettingsData.mutate,
        isValidating: newArgs.listFederatedSettingsData.isValidating,
      };
    }
  );
  return (
    <DiProvider
      use={[
        useListClusterDetailsDi,
        useListFederatedVirtualServicesDi,
        useListFederatedRouteTablesDi,
        useListFederatedUpstreamsDi,
        useListFederatedUpstreamGroupsDi,
        useListFederatedAuthConfigsDi,
        useListFederatedGatewaysDi,
        useListFederatedSettingsDi,
      ]}>
      <MemoryRouter>
        <AdminLanding />
      </MemoryRouter>
    </DiProvider>
  );
};

export const Primary = Template.bind({});
const clusterDetailsOne = Array.from({ length: 1 }).map(() => {
  return createClusterDetailsObj();
});
Primary.args = {
  listDetails: {
    data: clusterDetailsOne,
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  federatedData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  listFederatedRouteTablesData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  listFederatedUpstreamsData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  listFederatedUpstreamGroupsData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  listFederatedAuthConfigsData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  listFederatedGatewaysData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
  listFederatedSettingsData: {
    data: [],
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
} as Partial<typeof AdminLanding>;
