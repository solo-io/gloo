import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedUpstreamGroups } from '../../Components/Features/Admin/FederatedUpstreamGroups';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createFederatedUpstreamGroup } from 'stories/mocks/generators';

export default {
  title: 'Admin / FederatedUpstreamGroups',
  component: FederatedUpstreamGroups,
} as ComponentMeta<typeof FederatedUpstreamGroups>;

const Template: ComponentStory<typeof FederatedUpstreamGroups> = (
  args: any
) => {
  const useListFederatedUpstreamGroupsDi = injectable(
    Apis.useListFederatedUpstreamGroups,
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
    <DiProvider use={[useListFederatedUpstreamGroupsDi]}>
      <FederatedUpstreamGroups />
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createFederatedUpstreamGroup();
});

export const Primary = Template.bind({});
Primary.args = {
  data,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof FederatedUpstreamGroups>;
