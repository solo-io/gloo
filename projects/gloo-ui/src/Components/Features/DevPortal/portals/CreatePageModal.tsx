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
import useSWR, { trigger } from 'swr';
import { useParams } from 'react-router';
import { portalApi } from '../api';

interface InitialPageCreationValuesType {
  name: string;
  url: string;
  description: string;
  linkName: string;
  displayOnHomepage: boolean;
}

const validationSchema = yup.object().shape({
  name: yup.string().required('The name is required'),
  url: yup
    .string()
    .required('The URL is required')
    .matches(/^(?:\/)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$/, {
      message: 'Must start with a backslash and follow with valid characters',
      excludeEmptyString: true
    }),
  linkName: yup.string().required('A name for the navigation link is required')
});

interface CreatePageModalProps {
  onClose: () => any;
}

export const CreatePageModal = (props: CreatePageModalProps) => {
  const { portalname, portalnamespace } = useParams();
  const { data: portal, error: portalListError } = useSWR(
    !!portalname && !!portalnamespace
      ? ['getPortal', portalname, portalnamespace]
      : null,
    (key, name, namespace) => portalApi.getPortalWithAssets({ name, namespace })
  );

  const [errorMessage, setErrorMessage] = React.useState('');

  const initialValues: InitialPageCreationValuesType = {
    name: '',
    url: '',
    description: '',
    linkName: '',
    displayOnHomepage: false
  };

  const attemptCreate = async (values: InitialPageCreationValuesType) => {
    portalApi
      .createPortalPage(
        { name: portalname!, namespace: portalnamespace! },
        {
          name: values.name,
          path: values.url,
          description: values.description,
          navigationLinkName: values.linkName,
          displayOnHomepage: values.displayOnHomepage
        }
      )
      .then(portal => {
        trigger(['getPortal', portalname, portalnamespace]);

        props.onClose();
      })
      .catch(err => {
        setErrorMessage(err);

        setTimeout(() => {
          setErrorMessage('');
        }, 10000);
      });
  };

  return (
    <div className='relative flex flex-col pt-4'>
      <Formik<InitialPageCreationValuesType>
        initialValues={initialValues}
        onSubmit={attemptCreate}
        validationSchema={validationSchema}>
        {({ handleSubmit, values }) => (
          <>
            <div className='flex items-center text-lg font-medium'>
              Create a Portal Page
              <span className='text-blue-600'>
                <PortalPageIcon className='w-6 h-6 ml-2 fill-current ' />
              </span>
            </div>

            <div className='p-3 mt-3 text-gray-700 bg-gray-100 rounded-lg'>
              Lorem Ipsum
            </div>

            {!!errorMessage.length && (
              <div className='p-4 text-orange-600 bg-orange-200'>
                {errorMessage}
              </div>
            )}

            <div className='grid grid-cols-2 gap-4 mt-4'>
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
            <div className='flex mb-2'>
              <div>
                <SoloFormCheckbox name={'displayOnHomepage'} hideError={true} />
                <span className='ml-1 mr-4 font-normal'>
                  Display on Home Page
                </span>
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
