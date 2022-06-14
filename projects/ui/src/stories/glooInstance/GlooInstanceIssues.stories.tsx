import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceIssues } from '../../Components/Features/GlooInstance/GlooInstanceIssues';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { ObjectMeta } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';

export default {
  title: 'GlooInstance / GlooInstanceIssues',
  component: GlooInstanceIssues,
} as ComponentMeta<typeof GlooInstanceIssues>;

const instance = new GlooInstance();
const meta = new ObjectMeta();
instance.setMetadata(meta);

const Template: ComponentStory<typeof GlooInstanceIssues> = args => (
  // @ts-ignore
  <GlooInstanceIssues {...args} />
);

export const Primary = Template.bind({});
Primary.args = {
  glooInstance: instance.toObject(),
} as Partial<typeof GlooInstanceIssues>;
