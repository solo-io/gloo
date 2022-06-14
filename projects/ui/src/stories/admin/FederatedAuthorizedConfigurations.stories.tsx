import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedAuthorizedConfigurations } from '../../Components/Features/Admin/FederatedAuthorizedConfigurations';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedAuthorizedConfigurations',
  component: FederatedAuthorizedConfigurations,
} as ComponentMeta<typeof FederatedAuthorizedConfigurations>;

const Template: ComponentStory<
  typeof FederatedAuthorizedConfigurations
> = args => <FederatedAuthorizedConfigurations />;

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedAuthorizedConfigurations>;
