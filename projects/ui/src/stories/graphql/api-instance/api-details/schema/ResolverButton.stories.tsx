import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ResolverButton from '../../../../../Components/Features/Graphql/api-instance/api-details/schema/ResolverButton';
import { MemoryRouter } from 'react-router';

export default {
  title: 'Graphql / api-instance / api-details / schema / ResolverButton',
  component: ResolverButton,
} as ComponentMeta<typeof ResolverButton>;

const Template: ComponentStory<typeof ResolverButton> = args => (
  <ResolverButton {...args} />
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {} as Partial<typeof ResolverButton>;
