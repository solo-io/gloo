import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import {
  SoloFormRadio,
  SoloFormRadioOption,
} from '../../Components/Common/SoloFormComponents';
import { Formik } from 'formik';

const radioOptions: SoloFormRadioOption[] = [
  {
    displayValue: 'radioOneDisplayValue',
    value: 'radioOneValue',
    subHeader: 'subheader',
  },
  {
    displayValue: 'radioTwoDisplayValue',
    value: 'radioTwoValue',
  },
];

export default {
  title: `Form / ${SoloFormRadio.name}`,
  component: SoloFormRadio,
} as unknown as ComponentMeta<typeof SoloFormRadio>;

const Template: ComponentStory<typeof SoloFormRadio> = args => (
  <Formik
    initialValues={{
      value: 'any',
    }}
    onSubmit={(val: any) => {
      console.log('Value called', val);
    }}>
    <SoloFormRadio {...args} />
  </Formik>
);

export const Primary = Template.bind({});
Primary.args = {
  name: 'name',
  isUpdate: false,
  title: 'title',
  options: radioOptions,
  horizontal: true,
  titleAbove: true,
  onChange: (result: any) => {
    console.log('onChange called with result', result);
  },
};
