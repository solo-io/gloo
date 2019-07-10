import styled from '@emotion/styled/macro';
import { Divider } from 'antd';
import {
  SoloFormInput,
  SoloFormTypeahead
} from 'Components/Common/SoloFormField';
import { Field } from 'formik';
import { NamespacesContext } from 'GlooIApp';
import * as React from 'react';
import { AWS_REGIONS } from 'utils/upstreamHelpers';
import * as yup from 'yup';
import { InputContainer } from './CreateUpstreamForm';

const InputItem = styled.div`
  display: flex;
  flex-direction: column;
`;

const SectionHeader = styled.div`
  font-size: 18px;
  font-weight: 500;
  margin-top: 10px;
`;

const StyledDivider = styled(Divider)`
  margin: 12px 0;
`;

// TODO combine with main initial values
const initialValues = {
  region: '',
  secretRefNamespace: '',
  secretRefName: ''
};

interface Props {}

export const awsValidationSchema = yup.object().shape({
  region: yup.string(),
  secretRefNamespace: yup.string(),
  secretRefName: yup.string()
});

export const AwsUpstreamForm: React.FC<Props> = () => {
  const namespaces = React.useContext(NamespacesContext);

  const awsRegions = AWS_REGIONS.map(item => item.name);

  return (
    <div>
      <SectionHeader> AWS Upstream Settings</SectionHeader>
      <StyledDivider />
      <InputContainer>
        <InputItem>
          <Field
            name='region'
            title='Region'
            presetOptions={awsRegions}
            component={SoloFormTypeahead}
          />
        </InputItem>
        <InputItem>
          <Field
            name='secretRefNamespace'
            title='Secret Ref Namespace'
            presetOptions={namespaces}
            component={SoloFormTypeahead}
          />
        </InputItem>
        <InputItem>
          <Field
            name='secretRefName'
            title='Secret Ref Name'
            placeholder='Secret Ref Name'
            component={SoloFormInput}
          />
        </InputItem>
      </InputContainer>
    </div>
  );
};
