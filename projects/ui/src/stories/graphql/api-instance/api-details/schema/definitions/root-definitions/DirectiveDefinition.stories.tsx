import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import DirectiveDefinition from '../../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/root-definitions/DirectiveDefinition';
import { MemoryRouter } from 'react-router';

export default {
  title:
    'Graphql / api-instance / api-details / schema / definitions / root-definitions / DirectiveDefinition',
  component: DirectiveDefinition,
} as ComponentMeta<typeof DirectiveDefinition>;

const Template: ComponentStory<typeof DirectiveDefinition> = args => (
  <DirectiveDefinition {...args} />
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  isEditable: false,
  schema: {},
  node: {
    name: {
      value: 'some value',
    },
  },
} as Partial<typeof DirectiveDefinition>;
