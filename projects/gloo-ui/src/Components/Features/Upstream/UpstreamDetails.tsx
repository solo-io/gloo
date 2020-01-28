import { css } from '@emotion/core';
import styled from '@emotion/styled';
import { Button, Spin } from 'antd';
import { ReactComponent as UpstreamIcon } from 'assets/upstream-icon-circle.svg';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import {
  SoloAWSSecretsList,
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/Form/SoloFormField';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloButton } from 'Components/Common/SoloButton';
import { SoloTable } from 'Components/Common/SoloTable';
import { Form, Formik } from 'formik';
import { TransformationTemplate } from 'proto/gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation_pb';
import React, { useState } from 'react';
import { useHistory, useParams } from 'react-router';
import { upstreamAPI } from 'store/upstreams/api';
import { colors } from 'Styles';
import { SoloButtonCSS } from 'Styles/CommonEmotions/button';
import useSWR, { mutate } from 'swr';
import { AWS_REGIONS, getUpstreamType } from 'utils/helpers';
import { ConfigurationToggle } from '../VirtualService/Details/VirtualServiceDetails';

const UpstreamDetailsContainer = styled.div`
  display: grid;
  grid-template-rows: 35px 250px 50px 1fr;
  grid-template-areas: 'ch' 'c' 'fh' 'f';
`;

const UpstreamIconSectionCardSize = styled(UpstreamIcon)`
  width: 33px !important;
  max-height: none !important;
`;

const ConfigSection = styled.div`
  display: grid;
  grid-template-columns: 1fr 2fr;
`;
const YamlSection = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
`;
const SettingsSection = styled.div``;

const SecuritySection = styled.div``;
const FunctionsSection = styled.div`
  grid-area: functions;
`;

const ConfigSectionHeader = styled.div`
  grid-area: ch;
  display: flex;
  align-items: center;
  justify-content: space-between;
  grid-auto-rows: 1fr;
`;

export const UpstreamDetails = () => {
  const [showConfig, setShowConfig] = useState(false);
  let { upstreamnamespace, upstreamname } = useParams();
  const history = useHistory();
  const [isLoading, setIsLoading] = useState(false);
  const { data: currentUpstream, error: upstreamError, isValidating } = useSWR(
    !!upstreamname && !!upstreamnamespace
      ? ['getUpstream', upstreamname, upstreamnamespace]
      : null,
    (key: string, upstreamname: string, upstreamnamespace: string) =>
      upstreamAPI.getUpstream({
        ref: { name: upstreamname, namespace: upstreamnamespace }
      })
  );

  const awsRegions = AWS_REGIONS.map(item => item.name);

  React.useEffect(() => {
    if (
      isLoading &&
      !!currentUpstream?.upstream?.aws?.lambdaFunctionsList.length
    ) {
      setIsLoading(false);
    }
  }, [isValidating]);
  if (!currentUpstream) {
    return <div>Loading...</div>;
  }

  function formatRestFunctions(
    data: [string, TransformationTemplate.AsObject][]
  ) {
    let dataUsed = [];
    dataUsed = data.map(([name, transformationTemplate]) => {
      return {
        name,
        body: transformationTemplate?.body?.text,
        headers: transformationTemplate?.headersMap.length
      };
    });
    return dataUsed;
  }

  function handleSaveYamlEdit(editedYaml: string) {
    mutate(
      ['getUpstream', upstreamname, upstreamnamespace],
      upstreamAPI.updateUpstreamYaml({
        editedYamlData: {
          ref: {
            name: upstreamname!,
            namespace: upstreamnamespace!
          },
          editedYaml: editedYaml
        }
      })
    );
  }

  function handleUpdateUpstream(values: any) {
    setIsLoading(true);
    mutate(
      ['getUpstream', upstreamname, upstreamnamespace],
      upstreamAPI.updateUpstream({
        upstreamInput: {
          ...currentUpstream?.upstream!,
          aws: {
            ...currentUpstream?.upstream!.aws!,
            region: values.region!
          },
          metadata: {
            ...currentUpstream?.upstream?.metadata!
          }
        }
      }),
      false
    );
  }
  // update or attempt to update

  // TODO: update settings/ revert settings buttons
  return (
    <Formik
      onSubmit={handleUpdateUpstream}
      enableReinitialize
      initialValues={{
        region: currentUpstream.upstream?.aws?.region,
        secret: {
          name: currentUpstream.upstream?.aws?.secretRef?.name || '',
          namespace: currentUpstream.upstream?.aws?.secretRef?.namespace || ''
        },
        rootCert: currentUpstream?.upstream?.sslConfig?.sslFiles?.rootCa,
        tlsCert: currentUpstream?.upstream?.sslConfig?.sslFiles?.tlsCert,
        tlsPrivateKey: currentUpstream?.upstream?.sslConfig?.sslFiles?.tlsKey
      }}>
      {formikStuff => (
        <Form>
          {/* <pre>{JSON.stringify(formikStuff, null, 2)}</pre> */}
          <div>
            <Breadcrumb />
            <SectionCard
              cardName={currentUpstream?.upstream?.metadata?.name || ''}
              logoIcon={<UpstreamIconSectionCardSize />}
              health={currentUpstream?.upstream?.status?.state}
              headerSecondaryInformation={[
                {
                  title: 'namespace',
                  value: upstreamnamespace!
                },
                {
                  title: 'type',
                  value: getUpstreamType(currentUpstream?.upstream!)
                }
              ]}
              healthMessage={
                currentUpstream?.upstream?.status &&
                currentUpstream?.upstream?.status!.reason.length
                  ? currentUpstream?.upstream?.status!.reason
                  : 'Service Status'
              }
              onClose={() => history.push(`/upstreams/`)}>
              <>
                <div
                  css={css`
                    display: flex;
                    flex-direction: row;
                    justify-content: flex-end;
                  `}>
                  <ConfigurationToggle onClick={() => setShowConfig(s => !s)}>
                    {showConfig ? `Hide ` : `View `}
                    YAML Configuration
                  </ConfigurationToggle>

                  <FileDownloadLink
                    fileContent={currentUpstream?.raw?.content!}
                    fileName={currentUpstream?.raw?.fileName!}
                  />
                </div>
                <div>
                  {showConfig && (
                    <ConfigDisplayer
                      content={currentUpstream?.raw?.content || ''}
                      asEditor
                      yamlError={undefined}
                      saveEdits={handleSaveYamlEdit}
                    />
                  )}
                </div>
                <UpstreamDetailsContainer>
                  <ConfigSectionHeader>
                    <div
                      css={css`
                        font-size: 18px;
                        color: black;
                      `}>
                      Configuration
                    </div>
                  </ConfigSectionHeader>
                  {/* different depending on the upstream type   */}
                  <>
                    <div
                      css={css`
                        grid-area: c;
                        background: ${colors.januaryGrey};
                        padding: 16px;
                        display: grid;
                        grid-template-columns: 1fr 2fr;
                        grid-gap: 15px;
                      `}>
                      {currentUpstream?.upstream?.aws && (
                        <div>
                          <p
                            css={css`
                              font-size: 18px;
                              color: black;
                            `}>
                            AWS Upstream Settings
                          </p>
                          <div
                            css={css`
                              background: white;
                              padding: 5px;
                            `}>
                            <SoloFormTypeahead
                              hideError
                              testId='aws-region'
                              name='region'
                              title='Region'
                              defaultValue={
                                currentUpstream?.upstream?.aws?.region
                              }
                              presetOptions={awsRegions.map(region => {
                                return { value: region };
                              })}
                            />
                            <SoloAWSSecretsList
                              hideError
                              defaultValue={
                                currentUpstream?.upstream?.aws?.secretRef?.name
                              }
                              testId='aws-secret'
                              name='secret'
                              type='aws'
                            />
                          </div>
                        </div>
                      )}
                      <div>
                        <p
                          css={css`
                            font-size: 18px;
                            color: black;
                          `}>
                          Security
                        </p>
                        <div
                          css={css`
                            background: white;
                            display: grid;
                            grid-template-columns: 1fr 1fr;
                            grid-template-rows: 1fr 1fr;
                          `}>
                          <div
                            css={css`
                              padding: 5px;
                            `}>
                            <SoloFormInput
                              hideError
                              name='tlsCert'
                              title='TLS Certificate'
                            />
                          </div>
                          <div
                            css={css`
                              padding: 5px;
                            `}>
                            <SoloFormInput
                              hideError
                              name='tlsPrivateKey'
                              title='TLS Private Key'
                            />
                          </div>

                          <div
                            css={css`
                              padding: 5px;
                            `}>
                            <SoloFormInput
                              hideError
                              name='rootCert'
                              title='Root Certificate'
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                    {(!!currentUpstream?.upstream?.aws?.lambdaFunctionsList ||
                      !!currentUpstream?.upstream?.kube?.serviceSpec) && (
                      <div
                        css={css`
                          grid-area: fh;
                          color: black;
                          font-size: 18px;
                          align-self: center;
                        `}>
                        Functions
                      </div>
                    )}
                    {currentUpstream?.upstream?.aws?.lambdaFunctionsList && (
                      <>
                        <div
                          css={css`
                            grid-area: f;
                          `}>
                          <Spin
                            spinning={
                              isLoading ||
                              !currentUpstream?.upstream?.aws
                                ?.lambdaFunctionsList
                            }>
                            <SoloTable
                              dataSource={
                                currentUpstream?.upstream?.aws
                                  ?.lambdaFunctionsList
                              }
                              columns={[
                                {
                                  title: 'Name',
                                  dataIndex: 'lambdaFunctionName'
                                },
                                {
                                  title: 'Logical Name',
                                  dataIndex: 'logicalName'
                                },
                                {
                                  title: 'Qualifier',
                                  dataIndex: 'qualifier'
                                }
                              ]}
                            />
                          </Spin>
                        </div>
                      </>
                    )}

                    {currentUpstream?.upstream?.kube?.serviceSpec?.rest
                      ?.transformationsMap && (
                      <div
                        css={css`
                          grid-area: f;
                        `}>
                        <SoloTable
                          columns={[
                            {
                              title: 'Name',
                              dataIndex: 'name'
                            },
                            {
                              title: 'Body',
                              dataIndex: 'body'
                            },
                            {
                              title: 'Headers',
                              dataIndex: 'headers'
                            }
                          ]}
                          dataSource={formatRestFunctions(
                            currentUpstream?.upstream?.kube.serviceSpec.rest
                              .transformationsMap
                          )}
                        />
                      </div>
                    )}
                    <div
                      css={css`
                        display: grid;
                        grid-template-columns: 150px 150px;
                        grid-gap: 15px;
                        margin-top: 10px;
                      `}>
                      <Button
                        css={SoloButtonCSS}
                        type='primary'
                        htmlType='submit'>
                        Update Settings
                      </Button>
                      <SoloButton
                        css={css`
                          background-color: ${colors.juneGrey};
                          &:disabled {
                            background: ${colors.juneGrey};
                          }

                          &:hover {
                            background: ${colors.juneGrey};
                          }

                          &:focus,
                          &:active {
                            background: ${colors.juneGrey};
                          }
                        `}
                        onClick={() => {}}
                        text='Cancel'
                        // disabled={!editMode}
                      />
                    </div>
                  </>
                </UpstreamDetailsContainer>
              </>
            </SectionCard>
          </div>
        </Form>
      )}
    </Formik>
  );
};
