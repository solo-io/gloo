import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import SelectionSet from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/SelectionSet';
import { MemoryRouter } from 'react-router';
import { Kind } from 'graphql';

export default {
  title:
    'Graphql / api-instance / api-details / schema / definitions / SelectionSet',
  component: SelectionSet,
} as ComponentMeta<typeof SelectionSet>;

const Template: ComponentStory<typeof SelectionSet> = args => (
  <SelectionSet {...args} />
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
} as Partial<typeof SelectionSet>;
