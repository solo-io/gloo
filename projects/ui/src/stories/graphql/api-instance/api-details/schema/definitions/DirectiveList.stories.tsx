import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import DirectiveList from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/DirectiveList';
import { MemoryRouter } from 'react-router';

export default {
  title: `Graphql / api-instance / api-details / schema / definitions / ${DirectiveList.name}`,
  component: DirectiveList,
} as unknown as ComponentMeta<typeof DirectiveList>;

const Template: ComponentStory<typeof DirectiveList> = args => (
  <DirectiveList {...args} />
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  isEditable: true,
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
} as Partial<typeof DirectiveList>;
