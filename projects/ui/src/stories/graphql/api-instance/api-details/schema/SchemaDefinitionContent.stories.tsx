import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import SchemaDefinitionContent from '../../../../../Components/Features/Graphql/api-instance/api-details/schema/SchemaDefinitionContent';
// @ts-ignore
import { mockEnumDefinitions } from 'Components/Features/Graphql/api-instance/api-details/schema/mockData';

export default {
  title: `Graphql / api-instance / api-details / schema / ${SchemaDefinitionContent.name}`,
  component: SchemaDefinitionContent,
} as unknown as ComponentMeta<typeof SchemaDefinitionContent>;

const Template: ComponentStory<typeof SchemaDefinitionContent> = args => (
  <SchemaDefinitionContent {...args} />
);

export const Primary = Template.bind({});

const field = { ...mockEnumDefinitions[0] };
field.directives = [];
(field as any).definitions = [];

// @ts-ignore
Primary.args = {
  isEditable: false,
  node: field,
} as Partial<typeof SchemaDefinitionContent>;
