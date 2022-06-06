import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import NameAndReturnType from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/NameAndReturnType';
import { MemoryRouter } from 'react-router';
import { Kind } from 'graphql';

export default {
  title: `Graphql / api-instance / api-details / schema / definitions / ${NameAndReturnType.name}`,
  component: NameAndReturnType,
} as unknown as ComponentMeta<typeof NameAndReturnType>;

const Template: ComponentStory<typeof NameAndReturnType> = args => (
  <NameAndReturnType {...args} />
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
} as Partial<typeof NameAndReturnType>;
