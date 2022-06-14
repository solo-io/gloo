import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ResolverConfigSection } from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/ResolverConfigSection';
import { Formik } from 'formik';
import { expect, jest } from '@storybook/jest';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import * as Apis from 'API/hooks';
import { useAppSettings } from 'Components/Context/AppSettingsContext';
export default {
  title:
    'Graphql / api-instance / api-details / schema / resolver-wizard / ResolverConfigSection',
  component: ResolverConfigSection,
} as ComponentMeta<typeof ResolverConfigSection>;

const Template: ComponentStory<typeof ResolverConfigSection> = args => {
  const useGetConsoleOptionsDi = injectable(Apis.useGetConsoleOptions, () => {
    return {
      readonly: false,
      apiExplorerEnabled: true,
      errorMessage: '',
    };
  });
  const useAppSettingsDi = injectable(useAppSettings, () => {
    return {
      appSettings: {
        keyboardHandler: '',
      },
      onAppSettingsChange: jest.fn(),
    };
  });
  return (
    <DiProvider use={[useGetConsoleOptionsDi, useAppSettingsDi]}>
      <Formik
        onSubmit={() => {}}
        initialValues={{
          resolverType: {},
          upstream: '',
          resolverConfig: '',
          listOfResolvers: [],
          protoFile: '',
        }}>
        <ResolverConfigSection {...args} />
      </Formik>
    </DiProvider>
  );
};

export const Primary = Template.bind({});

Primary.args = {
  globalWarningMessage: '',
  onCancel: jest.fn(),
  formik: {
    values: {
      resolverConfig: '',
      resolverType: '',
    },
    isValid: true,
    errors: {},
  },
} as Partial<typeof ResolverConfigSection>;
