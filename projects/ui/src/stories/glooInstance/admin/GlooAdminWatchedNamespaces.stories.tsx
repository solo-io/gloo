import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminWatchedNamespaces } from '../../../Components/Features/GlooInstance/Admin/GlooAdminWatchNamespaces';

export default {
  title: `GlooInstance / Admin / ${GlooAdminWatchedNamespaces.name}`,
  component: GlooAdminWatchedNamespaces,
} as unknown as ComponentMeta<typeof GlooAdminWatchedNamespaces>;

const Template: ComponentStory<typeof GlooAdminWatchedNamespaces> = args => (
  // @ts-ignore
  <GlooAdminWatchedNamespaces {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminWatchedNamespaces>;
