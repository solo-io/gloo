import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import FieldList from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/FieldList';
import { MemoryRouter } from 'react-router';

export default {
  title: `Graphql / api-instance / api-details / schema / definitions / ${FieldList.name}`,
  component: FieldList,
} as unknown as ComponentMeta<typeof FieldList>;

const Template: ComponentStory<typeof FieldList> = args => (
  <FieldList {...args} />
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  isEditable: true,
  objectType: '',
  isRoot: false,
  node: {
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
} as Partial<typeof FieldList>;
