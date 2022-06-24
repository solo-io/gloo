import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedRouteTables } from '../../Components/Features/Admin/FederatedRouteTables';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createListFederatedRouteTables } from 'stories/mocks/generators';

export default {
  title: 'Admin / FederatedRouteTables',
  component: FederatedRouteTables,
} as ComponentMeta<typeof FederatedRouteTables>;

const Template: ComponentStory<typeof FederatedRouteTables> = (args: any) => {
  const useListFederatedRouteTablesDi = injectable(
    Apis.useListFederatedRouteTables,
    () => {
      return {
        data: args.tables,
        error: args.error,
        mutate: args.mutate,
        isValidating: args.isValidating,
      };
    }
  );
  return (
    <DiProvider use={[useListFederatedRouteTablesDi]}>
      <FederatedRouteTables />
    </DiProvider>
  );
};

const federatedRouteTablesList = Array.from({ length: 1 }).map(() => {
  return createListFederatedRouteTables();
});

export const Primary = Template.bind({});
Primary.args = {
  tables: federatedRouteTablesList,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof FederatedRouteTables>;
