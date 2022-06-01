import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import {
  SectionCard,
  SectionCardProps,
} from '../../Components/Common/SectionCard';

export default {
  title: `Card / ${SectionCard.name}`,
  component: SectionCard,
} as unknown as ComponentMeta<typeof SectionCard>;

const Template: ComponentStory<typeof SectionCard> = args => (
  <SectionCard {...args} />
);

export const Primary = Template.bind({});
Primary.args = {
  cardName: 'SectionCard',
} as Partial<SectionCardProps>;
