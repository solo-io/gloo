import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedAuthorizedConfigurations } from '../../Components/Features/Admin/FederatedAuthorizedConfigurations';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createFederatedAuthConfig } from 'stories/mocks/generators';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { FederatedAuthConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_enterprise_gloo_resources_pb';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedAuthorizedConfigurations',
  component: FederatedAuthorizedConfigurations,
} as ComponentMeta<typeof FederatedAuthorizedConfigurations>;

type FederatedAuthorizedConfigurationsStoryProps = {
  listFederatedAuthConfigsData: {
    data: FederatedAuthConfig.AsObject[];
    error?: ServiceError;
    isValidating: boolean;
    mutate: any;
  };
};

const Template: ComponentStory<typeof FederatedAuthorizedConfigurations> = (
  args: any
) => {
  const newArgs: FederatedAuthorizedConfigurationsStoryProps = { ...args };
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
  return (
    <DiProvider use={[useListFederatedAuthConfigsDi]}>
      <FederatedAuthorizedConfigurations />
    </DiProvider>
  );
};

const listFederatedAuthConfigsData = Array.from({ length: 2 }).map(() => {
  return createFederatedAuthConfig();
});

export const Primary = Template.bind({});
Primary.args = {
  listFederatedAuthConfigsData: {
    data: listFederatedAuthConfigsData,
    error: undefined,
    isValidating: false,
    mutate: jest.fn(),
  },
} as Partial<typeof FederatedAuthorizedConfigurations>;
