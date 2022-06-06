import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import VariableDefinitions from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/VariableDefinitions';
import { Kind } from 'graphql';

export default {
  title: `Graphql / api-instance / api-details / schema / definitions / ${VariableDefinitions.name}`,
  component: VariableDefinitions,
} as unknown as ComponentMeta<typeof VariableDefinitions>;

const Template: ComponentStory<typeof VariableDefinitions> = args => (
  <VariableDefinitions {...args} />
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  isEditable: true,
  objectType: '',
  canAddResolverThisLevel: false,
  isRoot: false,
  node: {
    kind: Kind.NAME,
    directives: [
      {
        name: {
          value: 'some foo',
        },
      },
      {
        name: {
          value: 'other bar',
        },
      },
    ],
  },
} as Partial<typeof VariableDefinitions>;
