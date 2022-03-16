import { graphqlConfigApi } from 'API/graphql';
import { useGetGraphqlApiDetails } from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import ErrorModal from 'Components/Common/ErrorModal';
import { SoloToggleSwitch } from 'Components/Common/SoloToggleSwitch';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect } from 'react';

const GraphqlApiIntrospectionToggle: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  // gets the schema from the api
  const { data: graphqlApi, error: graphqlApiError } =
    useGetGraphqlApiDetails(apiRef);

  // Updates schema
  const [attemptUpdateSchema, setAttemptUpdateSchema] = React.useState(false);
  const [introspectionEnabled, setIntrospectionEnabled] = React.useState(
    graphqlApi?.spec?.executableSchema?.executor?.local?.enableIntrospection ??
      false
  );
  const [errorMessage, setErrorMessage] = React.useState('');
  const [errorModal, setErrorModal] = React.useState(false);
  const updateApi = async () => {
    await graphqlConfigApi
      .updateGraphqlApi({
        graphqlApiRef: apiRef,
        spec: {
          executableSchema: {
            executor: {
              //@ts-ignore
              local: {
                enableIntrospection: introspectionEnabled,
              },
            },
          },
        },
      })
      .then(() => {
        setAttemptUpdateSchema(false);
      })
      .catch(err => {
        setErrorModal(true);
        setErrorMessage(err?.message ?? '');
      });
  };

  useEffect(() => {
    setIntrospectionEnabled(
      graphqlApi?.spec?.executableSchema?.executor?.local?.enableIntrospection!
    );
  }, [
    !!graphqlApi?.spec?.executableSchema?.executor,
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap?.length,
  ]);

  return (
    <div className='inline-block'>
      <div className='flex items-end justify-end'>
        <span className='text-lg text-gray-900'>Schema Introspection:</span>
        <div className={'ml-2 -mt-5'}>
          <SoloToggleSwitch
            checked={introspectionEnabled}
            onChange={() => {
              setAttemptUpdateSchema(true);
              setIntrospectionEnabled(!introspectionEnabled);
            }}
          />
        </div>
      </div>
      <ConfirmationModal
        visible={attemptUpdateSchema}
        confirmPrompt='update this schema'
        confirmButtonText='Update'
        goForIt={updateApi}
        cancel={() => {
          setAttemptUpdateSchema(false);
          setIntrospectionEnabled(
            graphqlApi?.spec?.executableSchema?.executor?.local
              ?.enableIntrospection ?? false
          );
        }}
        isNegative
      />
      <ErrorModal
        cancel={() => setErrorModal(false)}
        visible={errorModal}
        errorDescription={errorMessage}
        errorMessage={'Failure updating Graphql Schema'}
        isNegative={true}
      />
    </div>
  );
};

export default GraphqlApiIntrospectionToggle;
