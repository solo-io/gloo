import React from 'react';
import { Formik } from 'formik';
import { ReactComponent as PortalPageIcon } from 'assets/portal-page-icon.svg';
import {
  SoloFormInput,
  SoloFormTextarea,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import { SoloButtonStyledComponent } from 'Styles/CommonEmotions/button';
import * as yup from 'yup';

interface InitialPageCreationValuesType {
  name: string;
  url: string;
  description: string;
  linkName: string;
  useTopNav: boolean;
  useFooterNav: boolean;
}

const validationSchema = yup.object().shape({
  name: yup.string().required('The name is required'),
  url: yup
    .string()
    .required('The URL is required')
    .matches(/^(?:\/)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$/, {
      message: 'Organization may only include lowercase letters',
      excludeEmptyString: true
    }),
  linkName: yup.string().required('A name for the navigation link is required')
});

export const CreatePageModal = () => {
  const initialValues: InitialPageCreationValuesType = {
    name: '',
    url: '',
    description: '',
    linkName: '',
    useTopNav: true,
    useFooterNav: true
  };

  const attemptCreate = (values: InitialPageCreationValuesType) => {
    console.log(values);
  };

  return (
    <div className='relative flex flex-col pt-4'>
      <Formik<InitialPageCreationValuesType>
        initialValues={initialValues}
        onSubmit={attemptCreate}
        validationSchema={validationSchema}>
        {({ handleSubmit, values }) => (
          <>
            <div className='text-lg flex items-center font-medium'>
              Create a Portal Page
              <span className='text-blue-600'>
                <PortalPageIcon className=' ml-2 fill-current w-6 h-6' />
              </span>
            </div>

            <div className='rounded-lg bg-gray-100 p-3 mt-3 text-gray-700'>
              Lorem Ipsum
            </div>

            <div className='mt-4 grid grid-cols-2 gap-4'>
              <div className='mb-4'>
                <SoloFormInput
                  testId={`create-portal-page-name`}
                  name='name'
                  placeholder='Gettiing Started'
                  title='Page Name'
                  hideError={true}
                />
              </div>
              <div className='mb-4'>
                <SoloFormInput
                  testId={`create-portal-page-url`}
                  name='url'
                  placeholder='/getting-started'
                  title='Page URL'
                  hideError={true}
                />
              </div>
            </div>
            <div>
              <SoloFormTextarea
                testId={`create-portal-page-description`}
                name='description'
                placeholder='This is the meta description for the page'
                title='Page Description'
                hideError={true}
              />
            </div>
            <div className='mb-2'>
              <SoloFormInput
                testId={`create-portal-page-linkname`}
                name='linkName'
                placeholder='Getting Started'
                title='Navigation Link Name'
                hideError={true}
              />
            </div>
            <div className='mb-2 flex'>
              <div>
                <SoloFormCheckbox name={'useTopNav'} hideError={true} />
                <span className='font-normal ml-1 mr-4'>Top Nav</span>
              </div>
              <div>
                <SoloFormCheckbox name={'useFooterNav'} hideError={true} />
                <span className='font-normal ml-1 mr-4'>Footer Nav</span>
              </div>
            </div>

            <div className='flex justify-end'>
              <SoloButtonStyledComponent onClick={handleSubmit}>
                Create Page
              </SoloButtonStyledComponent>
            </div>
          </>
        )}
      </Formik>
    </div>
  );
};
