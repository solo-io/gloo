import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { UpstreamSection } from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/UpstreamSection';
import { Formik } from 'formik';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { mockEnumDefinitions } from 'Components/Features/Graphql/api-instance/api-details/schema/mockData';

export default {
  title: `Graphql / api-instance / api-details / schema / resolver-wizard / ${UpstreamSection.name}`,
  component: UpstreamSection,
} as unknown as ComponentMeta<typeof UpstreamSection>;

const Template: ComponentStory<typeof UpstreamSection> = args => (
  <Formik
    onSubmit={() => {}}
    initialValues={{
      resolverType: {},
      upstream: '',
      resolverConfig: '',
      listOfResolvers: [],
      protoFile: '',
    }}>
    <UpstreamSection {...args} />
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
} as Partial<typeof UpstreamSection>;
