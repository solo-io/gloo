import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GraphqlLanding } from '../../Components/Features/Graphql/GraphqlLanding';

export default {
  title: `Graphql / ${GraphqlLanding.name}`,
  component: GraphqlLanding,
} as unknown as ComponentMeta<typeof GraphqlLanding>;

const Template: ComponentStory<typeof GraphqlLanding> = args => (
  // @ts-ignore
  <GraphqlLanding {...args} />
);

export const Primary = Template.bind({});

Primary.args = {} as Partial<typeof GraphqlLanding>;
