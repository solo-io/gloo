import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedSettingsTable } from '../../Components/Features/Admin/FederatedSettingsTable';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createFederatedSettings } from 'stories/mocks/generators';

export default {
  title: 'Admin / FederatedSettingsTable',
  component: FederatedSettingsTable,
} as ComponentMeta<typeof FederatedSettingsTable>;

const Template: ComponentStory<typeof FederatedSettingsTable> = (args: any) => {
  const useListFederatedSettingsDi = injectable(
    Apis.useListFederatedSettings,
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
    <DiProvider use={[useListFederatedSettingsDi]}>
      <FederatedSettingsTable />
    </DiProvider>
  );
};

const federatedSettings = Array.from({ length: 1 }).map(() => {
  return createFederatedSettings();
});

export const Primary = Template.bind({});
Primary.args = {
  data: federatedSettings,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof FederatedSettingsTable>;
