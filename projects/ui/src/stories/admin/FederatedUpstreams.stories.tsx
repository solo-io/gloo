import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedUpstreams } from '../../Components/Features/Admin/FederatedUpstreams';

// TODO:  Add in mock from jest
export default {
  title: `Admin / ${FederatedUpstreams.name}`,
  component: FederatedUpstreams,
} as unknown as ComponentMeta<typeof FederatedUpstreams>;

const Template: ComponentStory<typeof FederatedUpstreams> = args => (
  // @ts-ignore
  <FederatedUpstreams {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedUpstreams>;
