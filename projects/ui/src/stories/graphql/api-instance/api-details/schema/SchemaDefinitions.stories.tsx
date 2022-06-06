import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import SchemaDefinitions from '../../../../../Components/Features/Graphql/api-instance/api-details/schema/SchemaDefinitions';
// @ts-ignore
import { mockEnumDefinitions } from 'Components/Features/Graphql/api-instance/api-details/schema/mockData';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export default {
  title: `Graphql / api-instance / api-details / schema / ${SchemaDefinitions.name}`,
  component: SchemaDefinitions,
} as unknown as ComponentMeta<typeof SchemaDefinitions>;

const Template: ComponentStory<typeof SchemaDefinitions> = args => (
  <SchemaDefinitions {...args} />
);

export const Primary = Template.bind({});

const field = { ...mockEnumDefinitions[0] };
field.directives = [];
(field as any).definitions = [];
const apiRef = new ClusterObjectRef();

// @ts-ignore
Primary.args = {
  isEditable: false,
  node: field,
  schema: field,
  apiRef: apiRef.toObject(),
} as Partial<typeof SchemaDefinitions>;
