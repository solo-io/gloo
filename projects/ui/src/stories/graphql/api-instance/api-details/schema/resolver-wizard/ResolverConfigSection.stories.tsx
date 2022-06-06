import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ResolverConfigSection } from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/ResolverConfigSection';
import { Formik } from 'formik';

export default {
  title: `Graphql / api-instance / api-details / schema / resolver-wizard / ${ResolverConfigSection.name}`,
  component: ResolverConfigSection,
} as unknown as ComponentMeta<typeof ResolverConfigSection>;

const Template: ComponentStory<typeof ResolverConfigSection> = args => (
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
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  globalWarningMessage: '',
  formik: {
    values: {
      resolverConfig: '',
      resolverType: '',
    },
    isValid: true,
    errors: {},
  },
} as Partial<typeof ResolverConfigSection>;
