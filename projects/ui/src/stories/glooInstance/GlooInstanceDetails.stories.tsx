import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstancesDetails } from '../../Components/Features/GlooInstance/GlooInstanceDetails';

export default {
  title: `GlooInstance / ${GlooInstancesDetails.name}`,
  component: GlooInstancesDetails,
} as unknown as ComponentMeta<typeof GlooInstancesDetails>;

const Template: ComponentStory<typeof GlooInstancesDetails> = args => (
  // @ts-ignore
  <GlooInstancesDetails {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstancesDetails>;
