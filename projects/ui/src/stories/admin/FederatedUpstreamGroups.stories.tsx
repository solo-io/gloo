import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedUpstreamGroups } from '../../Components/Features/Admin/FederatedUpstreamGroups';

// TODO:  Add in mock from jest
export default {
  title: `Admin / ${FederatedUpstreamGroups.name}`,
  component: FederatedUpstreamGroups,
} as unknown as ComponentMeta<typeof FederatedUpstreamGroups>;

const Template: ComponentStory<typeof FederatedUpstreamGroups> = args => (
  // @ts-ignore
  <FederatedUpstreamGroups {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedUpstreamGroups>;
