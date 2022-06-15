import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedGateways } from '../../Components/Features/Admin/FederatedGateways';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import {
  createFederatedAuthConfig,
  createFederatedGateway,
} from 'stories/mocks/generators';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { FederatedAuthConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';
import { FederatedGateway } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_gateway_resources_pb';

type FederatedGatewaysProps = {
  listFederatedGatewaysData: {
    data: FederatedGateway.AsObject[];
    error?: ServiceError;
    mutate: any;
    isValidating: boolean;
  };
};

export default {
  title: 'Admin / FederatedGateways',
  component: FederatedGateways,
} as ComponentMeta<typeof FederatedGateways>;

const Template: ComponentStory<typeof FederatedGateways> = (args: any) => {
  const newArgs: FederatedGatewaysProps = { ...args };
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
  return (
    <DiProvider use={[useListFederatedGatewaysDi]}>
      <FederatedGateways />
    </DiProvider>
  );
};

export const Primary = Template.bind({});

const data = Array.from({ length: 1 }).map(() => {
  return createFederatedGateway();
});
Primary.args = {
  listFederatedGatewaysData: {
    data,
    error: undefined,
    mutate: jest.fn(),
    isValidating: false,
  },
} as Partial<typeof FederatedGateways>;
