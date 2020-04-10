import styled from '@emotion/styled';

import { Loading } from 'Components/Common/DisplayOnly/Loading';
import { SoloButton } from 'Components/Common/SoloButton';
import { Container } from 'Components/Features/Admin/AdminLanding';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import * as React from 'react';
import { IAceEditorProps } from 'react-ace';
import { useHistory, useParams } from 'react-router';
import useSWR from 'swr';
import yaml from 'yaml';
import { apiDocApi } from '../api';
import { SwaggerHolder } from './SwaggerHolder';

const Editor: React.FC<IAceEditorProps> = props => {
  if (typeof window !== 'undefined') {
    const Ace = require('react-ace').default;
    require('ace-builds/src-noconflict/mode-json');
    require('ace-builds/src-noconflict/theme-github');

    return <Ace {...props} />;
  }

  return null;
};

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

  const onContentChange = (code: string) => {
    setSwaggerString(code);
    setSwaggerContent(yaml.parse(code));
  };

  if (!apiDoc) {
    return <Loading>Loading...</Loading>;
  }

  return (
    <div>
      <div>
        <div className='flex justify-between mb-4'>
          <div className='text-2xl'>API Editor</div>
          <div className='flex'>
            <div className='mr-8'></div>
            <SoloButton text='Publish Changes' onClick={handleSave} />
          </div>
        </div>
        {!!saveError && <div className='mb-4'>{saveError}</div>}
      </div>
      <Columns>
        <Editor
          value={swaggerString}
          mode='json'
          fontSize={14}
          theme='github'
          onChange={onContentChange}
          name='UNIQUE_ID_OF_DIV'
          maxLines={50}
          width='100%'
          editorProps={{ $blockScrolling: true }}
        />
        <div>{!!swaggerContent && <SwaggerHolder spec={swaggerContent} />}</div>
      </Columns>
    </div>
  );
};
