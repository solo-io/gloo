import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Clusters } from '../../Components/Features/Admin/Clusters';

export default {
  title: `Admin / ${Clusters.name}`,
  component: Clusters,
} as unknown as ComponentMeta<typeof Clusters>;

const Template: ComponentStory<typeof Clusters> = args => (
  // @ts-ignore
  <Clusters {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof Clusters>;
