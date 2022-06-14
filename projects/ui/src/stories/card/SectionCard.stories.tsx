import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import {
  SectionCard,
  SectionCardProps,
} from '../../Components/Common/SectionCard';

export default {
  title: 'Card / SectionCard',
  component: SectionCard,
} as ComponentMeta<typeof SectionCard>;

const Template: ComponentStory<typeof SectionCard> = args => (
  <SectionCard {...args} />
);

export const Primary = Template.bind({});
Primary.args = {
  cardName: SectionCard.name,
} as Partial<SectionCardProps>;
