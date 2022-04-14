import { graphqlConfigApi } from 'API/graphql';
import { useGetGraphqlApiDetails, useGetConsoleOptions } from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import ErrorModal from 'Components/Common/ErrorModal';
import { SoloToggleSwitch } from 'Components/Common/SoloToggleSwitch';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect } from 'react';
import { SoloInput } from 'Components/Common/SoloInput';
import * as Styles from './GraphqlApiPolicy.style';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { useParams } from 'react-router';

const GraphqlApiPolicyInputs: React.FC = () => {
  // gets the schema from the api
  const { graphqlApiName, graphqlApiNamespace, graphqlApiClusterName } =
    useParams();
  const apiRef = {
    name: graphqlApiName,
    namespace: graphqlApiNamespace,
    clusterName: graphqlApiClusterName,
  };

  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);

  // Updates schema
  const [attemptUpdateSchema, setAttemptUpdateSchema] = React.useState(false);
  const [maxQueryDepth, setMaxQueryDepth] = React.useState<number>(
    graphqlApi?.spec?.executableSchema?.executor?.local?.options?.maxDepth
      ?.value ?? 0
  );

  const { readonly } = useGetConsoleOptions();
  const [introspectionEnabled, setIntrospectionEnabled] = React.useState(
    (graphqlApi?.spec?.executableSchema?.executor?.local?.enableIntrospection &&
      !readonly) ??
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
                options: {
                  maxDepth: {
                    value: maxQueryDepth,
                  },
                },
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
      graphqlApi?.spec?.executableSchema?.executor?.local
        ?.enableIntrospection! && !readonly
    );
  }, [
    !!graphqlApi?.spec?.executableSchema?.executor,
    graphqlApi?.spec?.executableSchema?.executor?.local?.resolutionsMap?.length,
    readonly,
  ]);

  const updateMaxQuery = (e: React.ChangeEvent<HTMLInputElement>) => {
    const {
      target: { value },
    } = e;
    let updatedValue = Number(value);
    if (updatedValue < 0) {
      updatedValue = 0;
    }
    setMaxQueryDepth(updatedValue);
  };

  const submitAttempt = () => {
    setAttemptUpdateSchema(true);
  };

  if (readonly || !graphqlApi) {
    return null;
  }

  return (
    <div className='mt-5 ml-5 mb-5'>
      <Styles.ItemWrapper>
        <Styles.LabelWrapper className={'ml-2'}>
          <Styles.InputLabel className='text-lg text-gray-900'>
            Schema Introspection:
          </Styles.InputLabel>
          <Styles.DescriptionText className=''>
            Allow clients to explore your schema using introspection.
          </Styles.DescriptionText>
        </Styles.LabelWrapper>
        <Styles.SwitchContainer className=''>
          <SoloToggleSwitch
            checked={introspectionEnabled}
            onChange={() => {
              setIntrospectionEnabled(!introspectionEnabled);
            }}
          />
        </Styles.SwitchContainer>
      </Styles.ItemWrapper>

      <Styles.ItemWrapper className='mt-10'>
        <Styles.LabelWrapper className={'ml-2'}>
          <Styles.InputLabel className='text-lg text-gray-900'>
            Maximum Query Depth:
          </Styles.InputLabel>
          <Styles.DescriptionText className=''>
            Limits the level of nesting allowed in queries.
          </Styles.DescriptionText>
        </Styles.LabelWrapper>
        <Styles.NumericContainer className=''>
          <SoloInput
            value={maxQueryDepth as any}
            onChange={updateMaxQuery}
            type='number'
            min='0'
          />
        </Styles.NumericContainer>
      </Styles.ItemWrapper>
      <Styles.ButtonContainer className='ml-2'>
        <SoloButtonStyledComponent onClick={submitAttempt}>
          Update Policies
        </SoloButtonStyledComponent>
      </Styles.ButtonContainer>
      <ConfirmationModal
        visible={attemptUpdateSchema}
        confirmPrompt='update this schema'
        confirmButtonText='Update'
        goForIt={updateApi}
        cancel={() => {
          setAttemptUpdateSchema(false);
          setIntrospectionEnabled(
            (graphqlApi?.spec?.executableSchema?.executor?.local
              ?.enableIntrospection &&
              !readonly) ??
              false
          );
          updateApi().catch(err => {
            setErrorMessage(err.message);
          });
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

export default GraphqlApiPolicyInputs;
