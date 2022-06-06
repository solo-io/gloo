import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import Description from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/definitions/Description';
import { MemoryRouter } from 'react-router';

export default {
  title: `Graphql / api-instance / api-details / schema / definitions / ${Description.name}`,
  component: Description,
} as unknown as ComponentMeta<typeof Description>;

const Template: ComponentStory<typeof Description> = args => (
  <Description {...args} />
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  isRoot: true,
  node: { description: { value: 'foo' } },
} as Partial<typeof Description>;
