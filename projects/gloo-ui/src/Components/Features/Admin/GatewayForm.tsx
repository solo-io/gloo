import styled from '@emotion/styled';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { SuccessModal } from 'Components/Common/DisplayOnly/SuccessModal';
import {
  SoloFormCheckbox,
  SoloFormDurationEditor,
  SoloFormInput
} from 'Components/Common/Form/SoloFormField';
import { SoloButton } from 'Components/Common/SoloButton';
import { Formik, FormikErrors } from 'formik';
import { HttpConnectionManagerSettings } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/options/hcm/hcm_pb';
import {
  EditedResourceYaml,
  Raw
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb';
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppState } from 'store';
import { updateGatewayYaml } from 'store/gateway/actions';
import { colors, soloConstants } from 'Styles';
import * as yup from 'yup';
import { Gateway } from 'proto/github.com/solo-io/gloo/projects/gateway/api/v1/gateway_pb';
import * as github_com_solo_io_gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb
  from "../../../proto/github.com/solo-io/gloo/projects/gloo/api/v1/options/protocol_upgrade/protocol_upgrade_pb";

const GatewayFormContainer = styled.div`
  background: ${colors.januaryGrey};
  border: 1px solid ${colors.marchGrey};
  border-radius: ${soloConstants.smallRadius}px;
  padding: 13px;
  margin-top: 13px;
  color: ${colors.juneGrey};
  margin-bottom: 15px;
`;

const ExpandableSection = styled.div`
  max-height: ${(props: { isExpanded: boolean }) =>
    props.isExpanded ? '1000px' : '0px'};
  overflow: ${(props: { isExpanded: boolean }) =>
    props.isExpanded ? 'auto' : 'hidden'};
  transition: max-height ${soloConstants.transitionTime};
  color: ${colors.septemberGrey};
`;

const InnerSectionTitle = styled.div`
  color: ${colors.novemberGrey};
  font-size: 18px;
  line-height: 22px;
  margin: 13px 0;
`;

const InnerFormSectionContent = styled.div`
  background: white;
  border: 1px solid ${colors.marchGrey};
  border-radius: ${soloConstants.smallRadius}px;
  margin-bottom: 13px;
  padding: 13px 8px 0;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  grid-gap: 8px;
`;

const FormItem = styled.div`
  display: flex;
  flex-direction: column;
`;

const FormFooter = styled.div`
  grid-column: 2;
  display: flex;
  justify-content: flex-end;
`;

const ConfigurationSection = styled.div`
  margin-top: 15px;
  overflow: auto;
`;
const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
  margin-bottom: 5px;
`;

export type HttpConnectionManagerSettingsForm = HttpConnectionManagerSettings.AsObject;

let defaultHttpValues: HttpConnectionManagerSettingsForm = {
  skipXffAppend: (undefined as unknown) as boolean,
  maxRequestHeadersKb: { value: (undefined as unknown) as number },
  streamIdleTimeout: undefined,
  via: (undefined as unknown) as string,
  requestTimeout: undefined,
  idleTimeout: undefined,
  xffNumTrustedHops: (undefined as unknown) as number,
  drainTimeout: undefined,
  defaultHostForHttp10: (undefined as unknown) as string,
  useRemoteAddress: { value: (undefined as unknown) as boolean },
  delayedCloseTimeout: undefined,
  acceptHttp10: (undefined as unknown) as boolean,
  generateRequestId: { value: (undefined as unknown) as boolean },
  serverName: (undefined as unknown) as string,
  proxy100Continue: (undefined as unknown) as boolean,
  tracing: {
    requestHeadersForTagsList: (undefined as unknown) as string[],
    verbose: (undefined as unknown) as boolean
  },
  preserveExternalRequestId: (undefined as unknown) as boolean,
  setCurrentClientCertDetails: (undefined as unknown) as HttpConnectionManagerSettings.SetCurrentClientCertDetails.AsObject,
  upgradesList: (undefined as unknown) as Array<github_com_solo_io_gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.AsObject>,
  forwardClientCertDetails: (undefined as unknown) as HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap]
};

const validationSchema = yup.object().shape({
  skipXffAppend: yup.boolean(),
  maxRequestHeadersKb: yup.object().shape({ value: yup.number() }),
  streamIdleTimeout: yup
    .object()
    .shape({ nanos: yup.number(), seconds: yup.number() }),
  via: yup.string(),
  requestTimeout: yup
    .object()
    .shape({ nanos: yup.number(), seconds: yup.number() }),
  idleTimeout: yup
    .object()
    .shape({ nanos: yup.number(), seconds: yup.number() }),
  xffNumTrustedHops: yup.number(),
  drainTimeout: yup
    .object()
    .shape({ nanos: yup.number(), seconds: yup.number() }),
  defaultHostForHttp10: yup.string(),
  useRemoteAddress: yup.object().shape({ value: yup.boolean() }),
  delayedCloseTimeout: yup
    .object()
    .shape({ nanos: yup.number(), seconds: yup.number() }),
  acceptHttp10: yup.boolean(),
  generateRequestId: yup.object().shape({ value: yup.boolean() }),
  serverName: yup.string(),
  proxy100Continue: yup.boolean(),
  requestHeadersForTags: yup.string(),
  tracing: yup.object().shape({
    requestHeadersForTagsList: yup.string(),
    verbose: yup.boolean()
  })
});

interface FormProps {
  gatewayValues: Gateway.AsObject;
  gatewayConfiguration?: Raw.AsObject;
  doUpdate: (values: HttpConnectionManagerSettingsForm) => void;
  isExpanded: boolean;
}
export const GatewayForm = (props: FormProps) => {
  const [showSuccessModal, setShowSuccessModal] = React.useState(false);
  const [showConfiguration, setShowConfiguration] = React.useState(false);

  const dispatch = useDispatch();

  let initialValues: HttpConnectionManagerSettingsForm = {
    ...defaultHttpValues
  };

  if (
    props.gatewayValues.httpGateway &&
    props.gatewayValues.httpGateway.options &&
    props.gatewayValues.httpGateway.options.httpConnectionManagerSettings
  ) {
    let httpValues =
      props.gatewayValues.httpGateway.options.httpConnectionManagerSettings;

    initialValues.skipXffAppend = httpValues.skipXffAppend;
    initialValues.via = httpValues.via;
    initialValues.xffNumTrustedHops = httpValues.xffNumTrustedHops;
    if (httpValues.useRemoteAddress) {
      initialValues.useRemoteAddress!.value = httpValues.useRemoteAddress.value;
    }
    if (httpValues.generateRequestId) {
      initialValues.generateRequestId!.value =
        httpValues.generateRequestId.value;
    }
    initialValues.proxy100Continue = httpValues.proxy100Continue;
    initialValues.streamIdleTimeout = httpValues.streamIdleTimeout;
    initialValues.idleTimeout = httpValues.idleTimeout;
    if (httpValues.maxRequestHeadersKb) {
      initialValues.maxRequestHeadersKb!.value =
        httpValues.maxRequestHeadersKb.value;
    }
    initialValues.requestTimeout = httpValues.requestTimeout;
    initialValues.drainTimeout = httpValues.drainTimeout;
    initialValues.delayedCloseTimeout = httpValues.delayedCloseTimeout;
    initialValues.serverName = httpValues.serverName;
    initialValues.acceptHttp10 = httpValues.acceptHttp10;
    initialValues.defaultHostForHttp10 = httpValues.defaultHostForHttp10;
    if (httpValues.tracing) {
      initialValues.tracing!.requestHeadersForTagsList = httpValues.tracing!.requestHeadersForTagsList;
      initialValues.tracing!.verbose = httpValues.tracing.verbose;
    }
  }

  const invalid = (
    values: HttpConnectionManagerSettingsForm,
    errors: FormikErrors<HttpConnectionManagerSettingsForm>
  ) => {
    let isInvalid = false;

    return isInvalid;
  };
  const isDirty = (formIsDirty: boolean) => {
    return formIsDirty;
  };

  const saveYamlChange = (newYaml: string) => {
    let editedYamlData = {
      editedYaml: newYaml,
      ref: {
        name: props.gatewayValues.metadata!.name,
        namespace: props.gatewayValues.metadata!.namespace
      }
    } as EditedResourceYaml.AsObject;
    const updateGatewayYamlRequest = {
      editedYamlData
    };

    dispatch(updateGatewayYaml(updateGatewayYamlRequest));
  };
  const yamlError = useSelector(
    (state: AppState) => state.gateways.yamlParseError
  );

  return (
    <GatewayFormContainer>
      <SuccessModal
        visible={showSuccessModal}
        successMessage='Gateway updated successfully'
      />
      <ExpandableSection isExpanded={props.isExpanded}>
        <ConfigDisplayer
          content={props.gatewayConfiguration!.content}
          whiteBacked
          asEditor
          yamlError={yamlError}
          saveEdits={saveYamlChange}
        />

        {!!props.gatewayConfiguration && (
          <ConfigurationSection>
            <Link onClick={() => setShowConfiguration(s => !s)}>
              {showConfiguration ? 'Hide' : 'View'} Configuration Form
            </Link>
            {showConfiguration && (
              <>
                <div>
                  Below are gateway configuration settings you can update here.
                  For more information on these settings, please visit our{' '}
                  <a
                    href='https://gloo.solo.io/v1/github.com/solo-io/gloo/projects/gateway/api/v2/gateway.proto.sk/'
                    target='_blank'
                    rel='noopener noreferrer'>
                    hcm plugin documentation
                  </a>
                  .
                </div>
                <Formik
                  initialValues={initialValues}
                  validationSchema={validationSchema}
                  onSubmit={values => {
                    setShowSuccessModal(true);
                    props.doUpdate(values);
                    setShowSuccessModal(false);
                  }}>
                  {({
                    isSubmitting,
                    handleSubmit,
                    isValid,
                    errors,
                    dirty,
                    values
                  }) => {
                    return (
                      <>
                        <InnerSectionTitle>
                          Http Connection Manager Settings
                        </InnerSectionTitle>
                        <InnerFormSectionContent>
                          <FormItem>
                            <SoloFormInput
                              name={'maxRequestHeadersKb.value'}
                              title={'maxRequestHeadersKb'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormDurationEditor
                              value={values.streamIdleTimeout}
                              name={'streamIdleTimeout'}
                              title={'streamIdleTimeout'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormInput name={'via'} title={'via'} />
                          </FormItem>
                          <FormItem>
                            <SoloFormDurationEditor
                              value={values.requestTimeout}
                              name={'requestTimeout'}
                              title={'requestTimeout'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormDurationEditor
                              value={values.idleTimeout}
                              name={'idleTimeout'}
                              title={'idleTimeout'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormInput
                              type='number'
                              name={'xffNumTrustedHops'}
                              title={'xffNumTrustedHops'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormDurationEditor
                              value={values.drainTimeout}
                              name={'drainTimeout'}
                              title={'drainTimeout'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormInput
                              name={'defaultHostForHttp10'}
                              title={'defaultHostForHttp10'}
                            />
                          </FormItem>

                          <FormItem>
                            <SoloFormCheckbox
                              name={'useRemoteAddress.value'}
                              title={'useRemoteAddress'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormDurationEditor
                              value={values.delayedCloseTimeout}
                              name={'delayedCloseTimeout'}
                              title={'delayedCloseTimeout'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormCheckbox
                              name={'acceptHttp10'}
                              title={'acceptHttp10'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormCheckbox
                              name={'generateRequestId.value'}
                              title={'generateRequestId'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormInput
                              name={'serverName'}
                              title={'serverName'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormCheckbox
                              name={'proxy100Continue'}
                              title={'proxy100Continue'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormCheckbox
                              name={'skipXffAppend'}
                              title={'skipXffAppend'}
                            />
                          </FormItem>
                        </InnerFormSectionContent>
                        <InnerSectionTitle>Tracing Settings</InnerSectionTitle>
                        <InnerFormSectionContent>
                          <FormItem>
                            <SoloFormInput
                              name={'tracing.requestHeadersForTagsList'}
                              title={'requestHeadersForTags'}
                            />
                          </FormItem>
                          <FormItem>
                            <SoloFormCheckbox
                              name='tracing.verbose'
                              title={'verbose'}
                            />
                          </FormItem>
                        </InnerFormSectionContent>
                        <FormFooter>
                          <SoloButton
                            onClick={handleSubmit}
                            text='Update Configuration'
                            disabled={
                              isSubmitting ||
                              invalid(values, errors) ||
                              !isDirty(dirty)
                            }
                          />
                        </FormFooter>
                      </>
                    );
                  }}
                </Formik>
              </>
            )}
          </ConfigurationSection>
        )}
      </ExpandableSection>
    </GatewayFormContainer>
  );
};
