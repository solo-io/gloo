import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedVirtualServices } from '../../Components/Features/Admin/FederatedVirtualServices';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createFederatedVirtualService } from 'stories/mocks/generators';
// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedVirtualServices',
  component: FederatedVirtualServices,
} as ComponentMeta<typeof FederatedVirtualServices>;

const Template: ComponentStory<typeof FederatedVirtualServices> = (
  args: any
) => {
  const useListFederatedVirtualServicesDi = injectable(
    Apis.useListFederatedVirtualServices,
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
    <DiProvider use={[useListFederatedVirtualServicesDi]}>
      <FederatedVirtualServices />
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createFederatedVirtualService();
});

export const Primary = Template.bind({});
Primary.args = {
  data,
  mutate: jest.fn(),
  isValidating: false,
  error: undefined,
} as Partial<typeof FederatedVirtualServices>;
