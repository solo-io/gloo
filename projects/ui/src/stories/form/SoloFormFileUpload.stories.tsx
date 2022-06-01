import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import {
  SoloFormFileUpload,
  SoloFormFileUploadProps,
} from '../../Components/Common/SoloFormComponents';
import { Formik } from 'formik';

export default {
  title: `Form / ${SoloFormFileUpload.name}`,
  component: SoloFormFileUpload,
} as unknown as ComponentMeta<typeof SoloFormFileUpload>;

const Template: ComponentStory<typeof SoloFormFileUpload> = args => (
  <Formik
    initialValues={{
      name: {
        value: 'any',
      },
    }}
    onSubmit={(val: any) => {
      console.log('Value called', val);
    }}>
    <SoloFormFileUpload {...args} />
  </Formik>
);

export const Primary = Template.bind({});
Primary.args = {
  name: 'name',
  title: 'file upload item',
  isUpdate: false,
  horizontal: false,
  titleAbove: false,
  isDisabled: false,
  buttonLabel: 'button Label',
  fileType: 'application/json,application/x-yaml,text/*',
  'data-testid': 'testForm',
  // error: '',
} as Partial<SoloFormFileUploadProps>;
