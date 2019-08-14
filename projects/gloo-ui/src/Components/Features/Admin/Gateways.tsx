import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { Formik, FormikErrors } from 'formik';
import * as yup from 'yup';
import { useGetGatewayList } from 'Api/v2/useGatewayClientV2';
import { ReactComponent as GatewayLogo } from 'assets/gateway-icon.svg';
import { colors, soloConstants } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import { SoloFormInput } from 'Components/Common/Form/SoloFormField';
import { NamespacesContext } from 'GlooIApp';
import { GatewayDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/gateway_pb';
import {
  HttpGateway,
  Gateway
} from 'proto/github.com/solo-io/gloo/projects/gateway/api/v2/gateway_pb';

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

  const namespaces = React.useContext(NamespacesContext);
  const { data, loading, error, setNewVariables } = useGetGatewayList({
    namespaces: namespaces.namespacesList
  });
  const [allGateways, setAllGateways] = React.useState<
    GatewayDetails.AsObject[]
  >([]);
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

  const updateGateway = (values: HttpValuesType) => {};

  //console.log(allGateways);

  return (
    <React.Fragment>
      {allGateways.map((gateway, ind) => {
        return (
          <SectionCard
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
            ]}>
            <InsideHeader>
              <div>Configuration Settings</div> <div>gateway-ssl.yaml</div>
            </InsideHeader>
            <GatewayForm
              doUpdate={updateGateway}
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
  kipXffAppend: string;
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
  kipXffAppend: '',
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
  authLimitNumber: yup.string()
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
          hcm plugin documentation>.
        </a>
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
                    <FormItem>
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
                    <FormItem>
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
