import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedRateLimits } from '../../Components/Features/Admin/FederatedRateLimits';
import { createFederatedRateLimitConfig } from 'stories/mocks/generators';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import {
  createFederatedAuthConfig,
  createFederatedGateway,
} from 'stories/mocks/generators';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { FederatedRateLimitConfig } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/fed.rpc/v1/federated_ratelimit_resources_pb';

export default {
  title: 'Admin / FederatedRateLimits',
  component: FederatedRateLimits,
} as ComponentMeta<typeof FederatedRateLimits>;

type FederatedRateLimitsProps = {
  limitConfigs: FederatedRateLimitConfig.AsObject[];
  isValidating: boolean;
  mutate: any;
  error?: ServiceError;
};

const Template: ComponentStory<typeof FederatedRateLimits> = (args: any) => {
  const newArgs: FederatedRateLimitsProps = { ...args };
  const useListFederatedRateLimitsDi = injectable(
    Apis.useListFederatedRateLimits,
    () => {
      return {
        data: newArgs.limitConfigs,
        error: newArgs.error,
        mutate: newArgs.mutate,
        isValidating: newArgs.isValidating,
      };
    }
  );

  return (
    <DiProvider use={[useListFederatedRateLimitsDi]}>
      <FederatedRateLimits />
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createFederatedRateLimitConfig();
});

export const Primary = Template.bind({});
Primary.args = {
  limitConfigs: data,
  isValidating: false,
  mutate: jest.fn(),
} as Partial<typeof FederatedRateLimits>;
