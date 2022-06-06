import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ResolverTypeSection } from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/ResolverTypeSection';
import { Formik } from 'formik';

export default {
  title: `Graphql / api-instance / api-details / schema / resolver-wizard / ${ResolverTypeSection.name}`,
  component: ResolverTypeSection,
} as unknown as ComponentMeta<typeof ResolverTypeSection>;

const Template: ComponentStory<typeof ResolverTypeSection> = args => (
  <Formik
    onSubmit={() => {}}
    initialValues={{
      resolverType: {},
      upstream: '',
      resolverConfig: '',
      listOfResolvers: [],
      protoFile: '',
    }}>
    <ResolverTypeSection {...args} />
  </Formik>
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  setWarningMessage: () => {},
} as Partial<typeof ResolverTypeSection>;
