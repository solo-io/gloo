import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import VariableDefinitions from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/VariableDefinitions';
import { Kind } from 'graphql';
import { expect, jest } from '@storybook/jest';

export default {
  title:
    'Graphql / api-instance / api-details / schema / definitions / VariableDefinitions',
  component: VariableDefinitions,
} as ComponentMeta<typeof VariableDefinitions>;

const Template: ComponentStory<typeof VariableDefinitions> = args => (
  <VariableDefinitions {...args} />
);

export const Primary = Template.bind({});

Primary.args = {
  schema: {
    kind: Kind.DOCUMENT,
    definitions: [],
  },
  objectType: '',
  onReturnTypeClicked: jest.fn(),
  isEditable: true,
  canAddResolverThisLevel: false,
  isRoot: false,
  node: {
    kind: Kind.NAME,
    variableDefinitions: [],
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
