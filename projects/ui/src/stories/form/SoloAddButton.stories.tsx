import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { default as SoloAddButton } from '../../Components/Common/SoloAddButton';

// import { Button } from './Button';

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
  title: `Form / ${SoloAddButton.name}`,
  component: SoloAddButton,
  // More on argTypes: https://storybook.js.org/docs/react/api/argtypes
} as ComponentMeta<typeof SoloAddButton>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof SoloAddButton> = args => (
  <SoloAddButton {...args} />
);

export const Primary = Template.bind({});
// More on args: https://storybook.js.org/docs/react/writing-stories/args
Primary.args = {
  children: 'Here is Content',
  onClick: () => {
    console.log('Clicked button');
  },
};
