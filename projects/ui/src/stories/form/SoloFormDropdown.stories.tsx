import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import {
  SoloFormDropdown,
  SoloFormDropdownProps,
} from '../../Components/Common/SoloFormComponents';
import { Formik } from 'formik';
import { SelectValue } from 'antd/lib/select';

export default {
  title: 'Form / SoloFormDropdown',
  component: SoloFormDropdown,
} as ComponentMeta<typeof SoloFormDropdown>;

const options = [
  {
    key: 'one',
    value: 'one',
  },
  {
    key: 'two',
    value: 'two',
  },
];

const Template: ComponentStory<typeof SoloFormDropdown> = args => (
  <Formik
    initialValues={{
      name: {
        value: 'any',
      },
    }}
    onSubmit={(val: any) => {
      console.log('Value called', val);
    }}>
    <SoloFormDropdown {...args} />
  </Formik>
);

let selectedValue = 'two';

export const Primary = Template.bind({});
Primary.args = {
  name: 'name',
  searchable: true,
  title: 'storydropdown item',
  value: selectedValue,
  hideError: true,
  onChange: (newValue: SelectValue) => {
    console.log('selected value', selectedValue);
    selectedValue = newValue.toString();
  },
  disabled: false,
  error: '',
  options,
} as Partial<SoloFormDropdownProps>;
