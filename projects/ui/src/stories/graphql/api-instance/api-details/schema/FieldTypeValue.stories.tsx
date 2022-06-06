import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import FieldTypeValue from '../../../../../Components/Features/Graphql/api-instance/api-details/schema/FieldTypeValue';
import { MemoryRouter } from 'react-router';
import { mockEnumDefinitions } from 'Components/Features/Graphql/api-instance/api-details/schema/mockData';

export default {
  title: `Graphql / api-instance / api-details / schema / ${FieldTypeValue.name}`,
  component: FieldTypeValue,
} as unknown as ComponentMeta<typeof FieldTypeValue>;

const Template: ComponentStory<typeof FieldTypeValue> = args => (
  <MemoryRouter>
    <FieldTypeValue {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});

const field = { ...mockEnumDefinitions[0] };
field.directives = [];
(field as any).definitions = [];

// @ts-ignore
Primary.args = {
  schema: field,
  field,
} as Partial<typeof FieldTypeValue>;
