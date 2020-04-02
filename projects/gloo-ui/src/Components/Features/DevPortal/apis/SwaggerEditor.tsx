import * as React from 'react';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import yaml from 'yaml';
import { SoloTextarea } from 'Components/Common/SoloTextarea';
import { SwaggerHolder } from './SwaggerHolder';
import styled from '@emotion/styled';
import { SoloButton } from 'Components/Common/SoloButton';
import { useParams, useHistory } from 'react-router';
import useSWR from 'swr';
import { apiDocApi } from '../api';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { InverseButtonCSS } from 'Styles/CommonEmotions/button';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { Container } from 'Components/Features/Admin/AdminLanding';

const Columns = styled.div`
  display: grid;
  grid-template-columns: 50% 50%;
  grid-gap: 12px;
`;

const stringToContent = (swaggerSpec: string): object | null => {
  try {
    return JSON.parse(swaggerSpec);
  } catch {
    try {
      return yaml.parse(swaggerSpec);
    } catch {
      return null;
    }
  }
};

export const SwaggerEditor = () => {
  const { apiname, apinamespace } = useParams();
  const { data: apiDoc, error: apiDocError } = useSWR(
    !!apiname && !!apinamespace ? ['getApiDoc', apiname, apinamespace] : null,
    (key, name, namespace) =>
      apiDocApi.getApiDoc({ apidoc: { name, namespace }, withassets: true })
  );

  if (!apiDoc) {
    return <div>Loading...</div>;
  }

  return (
    <Container>
      <ErrorBoundary fallback={<div>There was an error in the API Editor</div>}>
        <StatefulSwaggerEditor
          apiDoc={apiDoc}
          apinamespace={apinamespace!}
          apiname={apiname!}
        />
      </ErrorBoundary>
    </Container>
  );
};

type EditorProps = {
  apiDoc: ApiDoc.AsObject;
  apinamespace: string;
  apiname: string;
};

const StatefulSwaggerEditor = ({
  apiDoc,
  apinamespace,
  apiname
}: EditorProps) => {
  const history = useHistory();
  const [swaggerString, setSwaggerString] = React.useState<string>(
    apiDoc?.spec?.dataSource?.inlineString || ''
  );
  const [swaggerContent, setSwaggerContent] = React.useState(
    stringToContent(apiDoc?.spec?.dataSource?.inlineString || '')
  );
  const [saveError, setSaveError] = React.useState('');

  const handlePreview = () => {
    setSwaggerContent(stringToContent(swaggerString));
  };

  const handleSave = () => {
    apiDocApi
      .updateApiDocContent(
        { namespace: apinamespace, name: apiname },
        swaggerString
      )
      .then(() => {
        history.push(`/dev-portal/apis/${apinamespace}/${apiname}`);
      })
      .catch((err: string) => {
        setSaveError(err);
        setTimeout(() => {
          setSaveError('');
        }, 3000);
      });
  };

  const handleSaveEdits = (str: string) => {
    setSwaggerString(str);
    setSwaggerContent(stringToContent(str));
  };

  if (!apiDoc) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <div>
        <div className='flex justify-between mb-4'>
          <div className='text-2xl'>API Editor</div>
          <div className='flex'>
            <div className='mr-8'>
              <SoloButton
                uniqueCss={InverseButtonCSS}
                text='Update Preview'
                onClick={handlePreview}
              />
            </div>
            <SoloButton text='Publish Changes' onClick={handleSave} />
          </div>
        </div>
        {!!saveError && <div className='mb-4'>{saveError}</div>}
      </div>
      <Columns>
        {/* TODO SUPPORT YAML */}
        {/* <ConfigDisplayer
          content={swaggerString}
          isJson={true}
          asEditor={true}
          saveEdits={handleSaveEdits}
          whiteBacked
        /> */}
        <SoloTextarea
          value={swaggerString}
          onChange={e => setSwaggerString(e.target.value)}
          rows={50}
        />
        <div>
          {!!swaggerContent && <SwaggerHolder swaggerJSON={swaggerContent} />}
        </div>
      </Columns>
    </div>
  );
};
