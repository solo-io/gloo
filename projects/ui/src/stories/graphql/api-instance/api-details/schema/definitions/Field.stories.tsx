import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import Field from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/Field';
import { MemoryRouter } from 'react-router';

export default {
  title: `Graphql / api-instance / api-details / schema / definitions / ${Field.name}`,
  component: Field,
} as unknown as ComponentMeta<typeof Field>;

const Template: ComponentStory<typeof Field> = args => <Field {...args} />;

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  isEditable: true,
  objectType: '',
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
} as Partial<typeof Field>;
