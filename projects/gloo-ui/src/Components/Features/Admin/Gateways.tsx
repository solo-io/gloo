import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';
import styled from '@emotion/styled/macro';
import { Formik, FormikErrors } from 'formik';
import * as yup from 'yup';
import { useGetGatewayList, useUpdateGateway } from 'Api/v2/useGatewayClientV2';
import { ReactComponent as GatewayLogo } from 'assets/gateway-icon.svg';
import { colors, soloConstants, healthConstants } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloButton } from 'Components/Common/SoloButton';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { YamlDisplayer } from 'Components/Common/DisplayOnly/YamlDisplayer';
import { SoloFormInput } from 'Components/Common/Form/SoloFormField';
import {
  GatewayDetails,
  UpdateGatewayRequest
} from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import {
  HttpGateway,
  Gateway
} from 'proto/github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb';
import { UpdateGatewayHttpData } from 'Api/v2/GatewayClient';
import { useSelector } from 'react-redux';
import { AppState } from 'store';

const InsideHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 5px;
  color: ${colors.novemberGrey};
`;

const GatewayLogoFullSize = styled(GatewayLogo)`
  width: 33px !important;
  max-height: none !important;
`;

const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
`;

interface Props {}

export const Gateways = (props: Props) => {
  const [gatewaysOpen, setGatewaysOpen] = React.useState<boolean[]>([]);

  const {
    config: { namespacesList }
  } = useSelector((state: AppState) => state);
  const {
    data: updateData,
    loading: updateLoading,
    error: updateError,
    setNewVariables: setNewUpdateVariables,
    dataObj: updateGatewayObj
  } = useUpdateGateway(null);
  const {
    data,
    loading,
    error,
    setNewVariables,
    dataObj: gatewayObj
  } = useGetGatewayList({
    namespaces: namespacesList
  });
  const [allGateways, setAllGateways] = React.useState<
    GatewayDetails.AsObject[]
  >([]);

  React.useEffect(() => {
    if (!!updateData) {
      setNewVariables({
        namespaces: namespacesList
      });
    }
  }, [updateLoading]);
  React.useEffect(() => {
    if (!!data) {
      const newGateways = data
        .toObject()
        .gatewayDetailsList.filter(gateway => !!gateway.gateway);
      setAllGateways(newGateways);
      setGatewaysOpen(Array.from({ length: newGateways.length }, () => false));
    }
  }, [loading]);

  if (!data || (!data && loading)) {
    return <div>Loading...</div>;
  }

  const toggleExpansion = (indexToggled: number) => {
    setGatewaysOpen(
      gatewaysOpen.map((isOpen, ind) => {
        if (ind !== indexToggled) {
          return false;
        }

        return !isOpen;
      })
    );
  };

  const updateGateway = (values: HttpValuesType, gatewayIndex: number) => {
    let updateGatewayData: UpdateGatewayHttpData = {
      acceptHttp10: values.acceptHttp10.toLowerCase() === 'true',
      defaultHostForHttp10: values.defaultHostForHttp10,
      delayedCloseTimeout: {
        seconds: 0,
        nanos: values.delayedCloseTimeout.length
          ? parseInt(values.delayedCloseTimeout)
          : 0
      },
      drainTimeout: {
        seconds: 0,
        nanos: values.drainTimeout.length ? parseInt(values.drainTimeout) : 0
      },
      generateRequestId: {
        value: values.generateRequestId.toLowerCase() === 'true'
      },
      idleTimeout: {
        seconds: 0,
        nanos: values.idleTimeout.length ? parseInt(values.idleTimeout) : 0
      },
      maxRequestHeadersKb: {
        value: values.maxRequestHeadersKb.length
          ? parseInt(values.maxRequestHeadersKb)
          : 0
      },
      proxy100Continue: values.proxy100Continue.toLowerCase() === 'true',
      requestTimeout: {
        seconds: 0,
        nanos: values.requestTimeout.length
          ? parseInt(values.requestTimeout)
          : 0
      },
      serverName: values.serverName,
      skipXffAppend: values.skipXffAppend.toLowerCase() === 'true',
      streamIdleTimeout: {
        seconds: 0,
        nanos: values.streamIdleTimeout.length
          ? parseInt(values.streamIdleTimeout)
          : 0
      },
      requestHeadersForTagsList: values.requestHeadersForTags.split(','),
      verbose: values.verbose.toLowerCase() === 'true',
      useRemoteAddress: {
        value: values.useRemoteAddress.toLowerCase() === 'true'
      },
      via: values.via,
      xffNumTrustedHops: values.xffNumTrustedHops.length
        ? parseInt(values.xffNumTrustedHops)
        : 0
    };

    setNewUpdateVariables({
      originalGateway: data.getGatewayDetailsList()[gatewayIndex].getGateway()!,
      updates: updateGatewayData
    });
  };

  return (
    <React.Fragment>
      {allGateways.map((gateway, ind) => {
        return (
          <SectionCard
            key={gateway.gateway!.gatewayProxyName + ind}
            cardName={gateway.gateway!.gatewayProxyName}
            logoIcon={<GatewayLogoFullSize />}
            headerSecondaryInformation={[
              {
                title: 'BindPort',
                value: gateway.gateway!.bindPort.toString()
              },
              {
                title: 'Namespace',
                value: gateway.gateway!.metadata!.namespace
              },
              { title: 'SSL', value: gateway.gateway!.ssl ? 'True' : 'False' }
            ]}
            health={
              gateway.gateway!.status
                ? gateway.gateway!.status!.state
                : healthConstants.Pending.value
            }
            healthMessage={'Gateway Status'}>
            <InsideHeader>
              <div>Configuration Settings</div>{' '}
              {!!gateway.raw && (
                <FileDownloadLink
                  fileName={gateway.raw.fileName}
                  fileContent={gateway.raw.content}
                />
              )}
            </InsideHeader>
            <GatewayForm
              doUpdate={(values: HttpValuesType) => updateGateway(values, ind)}
              gatewayValues={gateway.gateway!}
              isExpanded={gatewaysOpen[ind]}
            />
            <Link onClick={() => toggleExpansion(ind)}>
              {gatewaysOpen[ind] ? 'Hide' : 'View'} Settings
            </Link>
          </SectionCard>
        );
      })}
    </React.Fragment>
  );
};

const GatewayFormContainer = styled.div`
  background: ${colors.januaryGrey};
  border: 1px solid ${colors.marchGrey};
  border-radius: ${soloConstants.smallRadius}px;
  padding: 13px;
  margin-top: 13px;
  color: ${colors.juneGrey};
  margin-bottom: 15px;
`;

const ExpandableSection = styled<'div', { isExpanded: boolean }>('div')`
  max-height: ${props => (props.isExpanded ? '1000px' : '0px')};
  overflow: hidden;
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

interface HttpValuesType {
  skipXffAppend: string;
  maxRequestHeadersKb: string;
  streamIdleTimeout: string;
  via: string;
  requestTimeout: string;
  idleTimeout: string;
  xffNumTrustedHops: string;
  drainTimeout: string;
  defaultHostForHttp10: string;
  useRemoteAddress: string;
  delayedCloseTimeout: string;
  acceptHttp10: string;
  generateRequestId: string;
  serverName: string;
  proxy100Continue: string;
  requestHeadersForTags: string;
  verbose: string;
}

let defaultHttpValues: HttpValuesType = {
  skipXffAppend: '',
  maxRequestHeadersKb: '',
  streamIdleTimeout: '',
  via: '',
  requestTimeout: '',
  idleTimeout: '',
  xffNumTrustedHops: '',
  drainTimeout: '',
  defaultHostForHttp10: '',
  useRemoteAddress: '',
  delayedCloseTimeout: '',
  acceptHttp10: '',
  generateRequestId: '',
  serverName: '',
  proxy100Continue: '',
  requestHeadersForTags: '',
  verbose: ''
};

const connectionManagerList = Object.keys(defaultHttpValues).slice(0, -2);
const tracingList = Object.keys(defaultHttpValues).slice(-2);

const validationSchema = yup.object().shape({
  skipXffAppend: yup.string().oneOf(['true', 'True', 'false', 'False']),
  maxRequestHeadersKb: yup.number(),
  streamIdleTimeout: yup.number(),
  via: yup.string(),
  requestTimeout: yup.number(),
  idleTimeout: yup.number(),
  xffNumTrustedHops: yup.number(),
  drainTimeout: yup.number(),
  defaultHostForHttp10: yup.string(),
  useRemoteAddress: yup.string().oneOf(['true', 'True', 'false', 'False']),
  delayedCloseTimeout: yup.number(),
  acceptHttp10: yup.string().oneOf(['true', 'True', 'false', 'False']),
  generateRequestId: yup.string().oneOf(['true', 'True', 'false', 'False']),
  serverName: yup.string(),
  proxy100Continue: yup.string().oneOf(['true', 'True', 'false', 'False']),
  requestHeadersForTags: yup.string(),
  verbose: yup.string().oneOf(['true', 'True', 'false', 'False'])
});

interface FormProps {
  gatewayValues: Gateway.AsObject;
  doUpdate: (values: HttpValuesType) => void;
  isExpanded: boolean;
}
const GatewayForm = (props: FormProps) => {
  let httpValues = {};
  if (
    props.gatewayValues.httpGateway &&
    props.gatewayValues.httpGateway.plugins &&
    props.gatewayValues.httpGateway.plugins.httpConnectionManagerSettings
  ) {
    httpValues =
      props.gatewayValues.httpGateway.plugins.httpConnectionManagerSettings;
  }

  const initialValues: HttpValuesType = { ...defaultHttpValues, ...httpValues };

  const invalid = (
    values: HttpValuesType,
    errors: FormikErrors<HttpValuesType>
  ) => {
    let isInvalid = false;

    return isInvalid;
  };
  const isDirty = (formIsDirty: boolean) => {
    return formIsDirty;
  };

  return (
    <GatewayFormContainer>
      <div>
        Below are gateway configuration settings you can update here. For more
        information on these settings, please visit our{' '}
        <a
          href='https://gloo.solo.io/v1/github.com/solo-io/gloo/projects/gateway/api/v2/gateway.proto.sk/'
          target='_blank'>
          hcm plugin documentation
        </a>
        .
      </div>
      <ExpandableSection isExpanded={props.isExpanded}>
        <Formik
          initialValues={initialValues}
          validationSchema={validationSchema}
          onSubmit={props.doUpdate}>
          {({ isSubmitting, handleSubmit, isValid, errors, dirty, values }) => {
            return (
              <React.Fragment>
                <InnerSectionTitle>
                  Http Connection Manager Settings
                </InnerSectionTitle>
                <InnerFormSectionContent>
                  {connectionManagerList.map(fieldName => (
                    <FormItem key={fieldName}>
                      <SoloFormInput
                        key={fieldName}
                        name={fieldName}
                        title={fieldName}
                      />
                    </FormItem>
                  ))}
                </InnerFormSectionContent>
                <InnerSectionTitle>Tracing Settings</InnerSectionTitle>
                <InnerFormSectionContent>
                  {tracingList.map(fieldName => (
                    <FormItem key={fieldName}>
                      <SoloFormInput
                        key={fieldName}
                        name={fieldName}
                        title={fieldName}
                      />
                    </FormItem>
                  ))}
                </InnerFormSectionContent>
                <FormFooter>
                  <SoloButton
                    onClick={handleSubmit}
                    text='Update Configuration'
                    disabled={
                      isSubmitting || invalid(values, errors) || !isDirty(dirty)
                    }
                  />
                </FormFooter>
              </React.Fragment>
            );
          }}
        </Formik>
      </ExpandableSection>
    </GatewayFormContainer>
  );
};
