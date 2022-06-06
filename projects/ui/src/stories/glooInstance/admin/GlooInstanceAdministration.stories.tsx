import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceAdministration } from '../../../Components/Features/GlooInstance/Admin/GlooInstanceAdministration';

export default {
  title: `GlooInstance / Admin / ${GlooInstanceAdministration.name}`,
  component: GlooInstanceAdministration,
} as unknown as ComponentMeta<typeof GlooInstanceAdministration>;

const Template: ComponentStory<typeof GlooInstanceAdministration> = args => (
  // @ts-ignore
  <GlooInstanceAdministration {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstanceAdministration>;
