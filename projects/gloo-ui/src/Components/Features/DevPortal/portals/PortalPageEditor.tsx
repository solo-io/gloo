import React from 'react';
import { SoloInput } from 'Components/Common/SoloInput';
import { useParams, useHistory } from 'react-router';
import useSWR, { trigger } from 'swr';
import { portalApi } from '../api';
import { Formik } from 'formik';
import * as yup from 'yup';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { Breadcrumb } from 'antd';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as PortalPageIcon } from 'assets/portal-page-icon.svg';
import {
  SoloFormInput,
  SoloFormTextarea,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
  SoloNegativeButton
} from 'Styles/CommonEmotions/button';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { Portal } from 'proto/dev-portal/api/grpc/admin/portal_pb';
import ReactMde from 'react-mde';
import ReactDOM from 'react-dom';
import * as Showdown from 'showdown';
import { Loading } from 'Components/Common/DisplayOnly/Loading';

interface InitialPageEditingValuesType {
  name: string;
  path: string;
  description: string;
  navigationLinkName: string;
  useTopNav: boolean;
  useFooterNav: boolean;
  //published: boolean;
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

const converter = new Showdown.Converter({
  tables: true,
  simplifiedAutoLink: true,
  strikethrough: true,
  tasklists: true
});

// TODO ARTURO // TODO JOE :: This whole page has not been tested. I was not able to create a
//   portal page through the api, so I couldn't test this with real data. I did my best off
//   conjecture and theory, so... good luck! ::x
interface FormProps {
  portal: Portal.AsObject;
  togglePreviewState: (inPreview: boolean) => any;
}
const PortalPageEditorForm = ({ portal, togglePreviewState }: FormProps) => {
  const routerHistory = useHistory();
  const { portalname, portalnamespace, pagename } = useParams();

  const portalPage = portal?.spec?.staticPagesList.find(
    staticPage => staticPage.name === pagename
  )!;

  if (!portalPage) {
    return <div>{pagename} not found on this portal.</div>;
  }

  const [errorMessage, setErrorMessage] = React.useState('');
  const [attemptingToDelete, setAttemptingToDelete] = React.useState(false);
  const [markdownString, setMarkdownString] = React.useState(
    portalPage.content?.inlineString || ''
  );

  const initialValues: InitialPageEditingValuesType = {
    name: '',
    path: '',
    description: '',
    navigationLinkName: '',
    useTopNav: true,
    useFooterNav: true
  };

  const publishEdits = async (values: InitialPageEditingValuesType) => {
    portalApi
      .createPortalPage(
        { name: portalname!, namespace: portalnamespace! },
        {
          name: values.name,
          path: values.path,
          description: values.description,
          navigationLinkName: values.navigationLinkName,
          content: {
            inlineString: markdownString!,
            inlineBytes: '',
            fetchUrl: ''
          }
        }
      )
      .then((portalResponse: any) => {
        trigger([
          'getPortal',
          portal.metadata?.name,
          portal.metadata?.namespace
        ]);

        routerHistory.push(
          routerHistory.location.pathname.split('/page-editor/')[0]
        );
      })
      .catch(err => {
        setErrorMessage(err);

        setTimeout(() => {
          setErrorMessage('');
        }, 10000);
      });
  };

  const attemptDelete = (name: string) => {
    setAttemptingToDelete(true);
  };
  const cancelDelete = () => {
    setAttemptingToDelete(false);
  };
  const finishDelete = () => {
    portalApi
      .deletePortalPage(
        { name: portalname!, namespace: portalnamespace! },
        pagename!
      )
      .then(portal => {
        routerHistory.push(
          routerHistory.location.pathname.split('/page-editor/')[0]
        );
      });
    /*.catch(err => {
        setErrorMessage(err);

        setTimeout(() => {
          setErrorMessage('');
        }, 10000);
      });*/
  };

  return (
    <Formik<InitialPageEditingValuesType>
      initialValues={initialValues}
      onSubmit={publishEdits}
      validationSchema={validationSchema}>
      {({ handleSubmit, values, isValid, dirty }) => (
        <>
          {!!errorMessage.length && (
            <div className='p-4 text-orange-600 bg-orange-200'>
              {errorMessage}
            </div>
          )}

          <div className='grid grid-cols-pageEditor gap-4'>
            <div className='p-3 bg-gray-100 rounded-lg border-gray-400 border'>
              <div className='mb-4'>
                <SoloFormInput
                  testId={`create-portal-page-name`}
                  name='name'
                  placeholder='Gettiing Started'
                  title='Page Name'
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
              <div className='mb-4'>
                <SoloFormInput
                  testId={`create-portal-page-url`}
                  name='url'
                  placeholder='/getting-started'
                  title='Page URL'
                  hideError={true}
                />
              </div>
              <div className='flex mb-4'>
                <div>
                  <SoloFormCheckbox name={'useTopNav'} hideError={true} />
                  <span className='ml-1 mr-4 font-normal'>Top Navigation</span>
                </div>
                <div>
                  <SoloFormCheckbox name={'useFooterNav'} hideError={true} />
                  <span className='ml-1 mr-4 font-normal'>
                    Footer Navigation
                  </span>
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
            </div>
            <div>
              {/*
              // TOODO ARTURO // TODO JOE 
               ReactMde has a 'tab' which lets you switch between Preview and Editing.
               https://github.com/andrerpena/react-mde
               https://codesandbox.io/s/react-mde-latest-bm6p3

               use styled emotions to wrap this and then look for the class on the tabs ('.mde-tabs') and display: none
               then just tell it the selected tabs with a local state change based on any other buttons you want to provide
              */}
              <ReactMde
                value={markdownString}
                onChange={setMarkdownString}
                generateMarkdownPreview={markdown =>
                  Promise.resolve(converter.makeHtml(markdown))
                }
              />
            </div>
          </div>

          <div className='flex justify-between mt-6'>
            <div>
              <SoloButtonStyledComponent
                className='mr-4'
                onClick={togglePreviewState(true)}>
                Preview Changes
              </SoloButtonStyledComponent>
              <SoloButtonStyledComponent
                className='mr-4'
                green={true}
                onClick={handleSubmit}
                disable={isValid && !dirty && !!markdownString?.length}>
                Publish Changes
              </SoloButtonStyledComponent>
              <SoloCancelButton
                onClick={() =>
                  routerHistory.push(
                    routerHistory.location.pathname.split('/page-editor/')[0]
                  )
                }>
                Cancel
              </SoloCancelButton>
            </div>
            <div>
              <SoloNegativeButton onClick={attemptDelete}>
                Delete Page
              </SoloNegativeButton>
            </div>
          </div>

          <ConfirmationModal
            visible={!!attemptingToDelete}
            confirmationTopic={`delete ${pagename}`}
            confirmText='Delete'
            goForIt={finishDelete}
            cancel={cancelDelete}
            isNegative={true}
          />
        </>
      )}
    </Formik>
  );
};

export const PortalPageEditor = () => {
  const routerHistory = useHistory();
  const { portalname, portalnamespace, pagename } = useParams();

  const { data: portal, error: portalListError } = useSWR(
    !!portalname && !!portalnamespace
      ? ['getPortal', portalname, portalnamespace]
      : null,
    (key, name, namespace) => portalApi.getPortalWithAssets({ name, namespace })
  );

  const [inPreviewState, setInPreviewState] = React.useState(false);

  if (!portal) {
    return (
      <div>
        <Breadcrumb />
        <SectionCard cardName=''>
          <Loading center={true} />
        </SectionCard>
      </div>
    );
  }

  return (
    <ErrorBoundary
      fallback={
        <div>There was an error with the getting to the page editor</div>
      }>
      <div>
        <Breadcrumb />
        <SectionCard
          cardName={pagename!}
          logoIcon={
            <span className='text-blue-500'>
              <PortalPageIcon className='fill-current w-6 h-6' />
            </span>
          }
          headerSecondaryInformation={[
            {
              title: 'Modified',
              value: 'WE NEED THIS'
            }
          ]}
          onClose={() =>
            routerHistory.push(
              routerHistory.location.pathname.split('/page-editor/')[0]
            )
          }>
          <PortalPageEditorForm
            portal={portal!}
            togglePreviewState={setInPreviewState}
          />
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
