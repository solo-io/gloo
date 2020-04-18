import React from 'react';
import { SoloInput } from 'Components/Common/SoloInput';
import styled from '@emotion/styled';
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
import { Portal } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/portal_pb';
import ReactDOM from 'react-dom';
import ReactMarkdown from 'react-markdown';
import * as Showdown from 'showdown';
import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { SoloTextarea, Label } from 'Components/Common/SoloTextarea';

interface InitialPageEditingValuesType {
  name: string;
  path: string;
  description: string;
  navigationLinkName: string;
  displayOnHomepage: boolean;
}

const validationSchema = yup.object().shape({
  name: yup.string().required('The name is required'),
  path: yup
    .string()
    .required('The URL is required')
    .matches(/^(?:\/)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$/, {
      message: 'Must start with a backslash and follow with valid characters',
      excludeEmptyString: true
    }),
  navigationLinkName: yup
    .string()
    .required('A name for the navigation link is required')
});

interface FormProps {
  portal: Portal.AsObject;
}
const PortalPageEditorForm = ({ portal }: FormProps) => {
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
  const [inPreviewState, setInPreviewState] = React.useState(false);
  const [markdownString, setMarkdownString] = React.useState(
    (!!portalPage.content?.inlineString &&
      portalPage.content?.inlineString.toString()) ||
      ''
  );

  const initialValues: InitialPageEditingValuesType = {
    name: portalPage.name,
    path: portalPage.path,
    description: portalPage.description,
    navigationLinkName: portalPage.navigationLinkName,
    displayOnHomepage: portalPage.displayOnHomepage
  };

  const publishEdits = async (values: InitialPageEditingValuesType) => {
    portalApi
      .updatePortalPage(
        { name: portalname!, namespace: portalnamespace! },
        {
          name: values.name,
          path: values.path,
          description: values.description,
          navigationLinkName: values.navigationLinkName,
          displayOnHomepage: values.displayOnHomepage,
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
                  disabled='true'
                />
              </div>
              <div className='mb-2'>
                <SoloFormInput
                  testId={`create-portal-page-linkname`}
                  name='navigationLinkName'
                  placeholder='Getting Started'
                  title='Navigation Link Name'
                  hideError={true}
                />
              </div>
              <div className='mb-4'>
                <SoloFormInput
                  testId={`create-portal-page-url`}
                  name='path'
                  placeholder='/getting-started'
                  title='Page URL'
                  hideError={true}
                />
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
              <div>
                <SoloFormCheckbox
                  testId={`create-portal-page-homepage`}
                  name={'displayOnHomepage'}
                  hideError={true}
                />
                <span className='ml-1 mr-4 font-normal'>
                  Display on Home Page
                </span>
              </div>
            </div>
            <div>
              <div>
                {inPreviewState && <ReactMarkdown source={markdownString} />}
                {!inPreviewState && (
                  <SoloTextarea
                    placeholder='Markdown-formatted page content'
                    value={markdownString}
                    rows={25}
                    onChange={e => setMarkdownString(e.target.value)}
                  />
                )}
              </div>
            </div>
          </div>

          <div className='flex justify-between mt-6'>
            <div>
              <SoloButtonStyledComponent
                className='mr-4'
                onClick={() => setInPreviewState(!inPreviewState)}>
                {inPreviewState ? 'Edit' : 'Preview Changes'}
              </SoloButtonStyledComponent>
              <SoloButtonStyledComponent
                className='mr-4'
                green='true'
                onClick={handleSubmit}
                disable={(
                  isValid &&
                  !dirty &&
                  !!markdownString?.length
                ).toString()}>
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
          onClose={() =>
            routerHistory.push(
              routerHistory.location.pathname.split('/page-editor/')[0]
            )
          }>
          <PortalPageEditorForm portal={portal!} />
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
