import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminGatewaysBox } from '../../../Components/Features/GlooInstance/Admin/GlooInstanceAdminBoxSummary';
import { MemoryRouter } from 'react-router';

export default {
  title: `GlooInstance / Admin / ${GlooAdminGatewaysBox.name}`,
  component: GlooAdminGatewaysBox,
} as unknown as ComponentMeta<typeof GlooAdminGatewaysBox>;

const Template: ComponentStory<typeof GlooAdminGatewaysBox> = args => (
  // @ts-ignore
  <MemoryRouter>
    <GlooAdminGatewaysBox {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminGatewaysBox>;
