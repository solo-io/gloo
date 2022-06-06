import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GraphqlOverview } from '../../Components/Features/Graphql/GraphqlOverview';
import { MemoryRouter } from 'react-router';

export default {
  title: `Graphql / ${GraphqlOverview.name}`,
  component: GraphqlOverview,
} as unknown as ComponentMeta<typeof GraphqlOverview>;

const Template: ComponentStory<typeof GraphqlOverview> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <GraphqlOverview {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});

Primary.args = {} as Partial<typeof GraphqlOverview>;
