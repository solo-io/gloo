import { Table } from 'antd';
import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloFormInput,
  SoloFormSecretRefInput
} from 'Components/Common/Form/SoloFormField';
import {
  Footer,
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { SoloButton } from 'Components/Common/SoloButton';
import {
  Field,
  FieldArray,
  FieldArrayRenderProps,
  Formik,
  FormikProps
} from 'formik';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import * as React from 'react';
import { AZURE_AUTH_LEVELS } from 'utils/azureHelpers';
import * as yup from 'yup';
import { UpstreamSpec as AzureUpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { SoloFormDropdown } from 'Components/Common/Form/SoloFormField';
import styled from '@emotion/styled/macro';

const StyledInputRow = styled(InputRow)`
  justify-content: space-around;
`;

const IconContainer = styled.div`
  display: flex;
`;
export interface AzureValuesType {
  azureFunctionAppName: string;
  azureSecretRef: ResourceRef.AsObject;
}

export const azureInitialValues: AzureValuesType = {
  azureFunctionAppName: '',
  azureSecretRef: {
    name: '',
    namespace: 'gloo-system'
  }
};

interface Props {}

export const azureValidationSchema = yup.object().shape({
  azureFunctionAppName: yup.string(),
  azureSecretRefNamespace: yup.string(),
  azureSecretRefName: yup.string()
});

export const AzureUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Azure Upstream Settings'>
      <InputRow>
        <SoloFormInput
          name='azureFunctionAppName'
          title='Function App Name'
          placeholder='Function App Name'
        />
        <SoloFormSecretRefInput name='azureSecretRef' type='azure' />
      </InputRow>
    </SoloFormTemplate>
  );
};
