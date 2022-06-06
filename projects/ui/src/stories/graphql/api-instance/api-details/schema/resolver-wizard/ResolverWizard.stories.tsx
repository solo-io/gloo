import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ResolverWizard } from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/ResolverWizard';
import { Formik } from 'formik';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { mockEnumDefinitions } from 'Components/Features/Graphql/api-instance/api-details/schema/mockData';

export default {
  title: `Graphql / api-instance / api-details / schema / resolver-wizard / ${ResolverWizard.name}`,
  component: ResolverWizard,
} as unknown as ComponentMeta<typeof ResolverWizard>;

const Template: ComponentStory<typeof ResolverWizard> = args => (
  <Formik
    onSubmit={() => {}}
    initialValues={{
      resolverType: {},
      upstream: '',
      resolverConfig: '',
      listOfResolvers: [],
      protoFile: '',
    }}>
    <ResolverWizard {...args} />
  </Formik>
);

export const Primary = Template.bind({});

const apiRef = new ClusterObjectRef();

const field = { ...mockEnumDefinitions[0] };
field.directives = [];
(field as any).definitions = [];

// @ts-ignore
Primary.args = {
  apiRef: apiRef.toObject(),
  field,
  objectType: '',
} as Partial<typeof ResolverWizard>;
