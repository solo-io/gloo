import { graphqlConfigApi } from 'API/graphql';
import {
  useGetConsoleOptions,
  useGetGraphqlApiDetails,
  usePageApiRef,
} from 'API/hooks';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloToggleSwitch } from 'Components/Common/SoloToggleSwitch';
import { useConfirm } from 'Components/Context/ConfirmModalContext';
import React, { useEffect } from 'react';
import toast from 'react-hot-toast';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { hotToastError } from 'utils/hooks';
import * as Styles from './GraphqlApiPolicy.style';

const GraphqlApiPolicyInputs: React.FC = () => {
  // gets the schema from the api
  const apiRef = usePageApiRef();
  const confirm = useConfirm();

  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);

  // Updates schema
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
  const updateApi = async () => {
    await graphqlConfigApi.updateGraphqlApi({
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
        <SoloButtonStyledComponent
          onClick={() =>
            confirm({
              confirmPrompt: 'update this schema',
              confirmButtonText: 'Update',
              isNegative: true,
            }).then(() =>
              toast
                .promise(updateApi(), {
                  loading: 'Updating API...',
                  success: 'API updated!',
                  error: hotToastError,
                })
                .catch(() => {
                  setIntrospectionEnabled(
                    (graphqlApi?.spec?.executableSchema?.executor?.local
                      ?.enableIntrospection &&
                      !readonly) ??
                      false
                  );
                  updateApi();
                })
            )
          }>
          Update Policies
        </SoloButtonStyledComponent>
      </Styles.ButtonContainer>
    </div>
  );
};

export default GraphqlApiPolicyInputs;
