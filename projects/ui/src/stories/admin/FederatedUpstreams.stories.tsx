import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedUpstreams } from '../../Components/Features/Admin/FederatedUpstreams';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createFederatedUpstream } from 'stories/mocks/generators';

export default {
  title: 'Admin / FederatedUpstreams',
  component: FederatedUpstreams,
} as ComponentMeta<typeof FederatedUpstreams>;

const Template: ComponentStory<typeof FederatedUpstreams> = (args: any) => {
  const useListFederatedUpstreamsDi = injectable(
    Apis.useListFederatedUpstreams,
    () => {
      return {
        data: args.data,
        error: args.error,
        mutate: args.mutate,
        isValidating: args.isValidating,
      };
    }
  );
  return (
    <DiProvider use={[useListFederatedUpstreamsDi]}>
      <FederatedUpstreams />
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createFederatedUpstream();
});

export const Primary = Template.bind({});
Primary.args = {
  data,
  mutate: jest.fn(),
  isValdating: false,
  error: undefined,
} as Partial<typeof FederatedUpstreams>;
