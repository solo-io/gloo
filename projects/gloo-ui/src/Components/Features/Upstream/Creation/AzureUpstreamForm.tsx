import { ReactComponent as CloseX } from 'assets/close-x.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloFormInput,
  SoloSecretRefInput
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate,
  Footer
} from 'Components/Common/Form/SoloFormTemplate';
import {
  Field,
  FieldArray,
  FieldArrayRenderProps,
  FormikProps,
  Formik
} from 'formik';
import * as React from 'react';
import { AZURE_AUTH_LEVELS } from 'utils/azureHelpers';
import * as yup from 'yup';
import { UpstreamSpec as AzureUpstreamSpec } from '../../../../proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/azure/azure_pb';
import { SoloFormDropdown } from '../../../Common/Form/SoloFormField';
import { initialValues } from './CreateUpstreamForm';
import { SoloButton } from 'Components/Common/SoloButton';
import { ResourceRef } from 'proto/github.com/solo-io/solo-kit/api/v1/ref_pb';
import { Table } from 'antd';

interface AzureValuesType {
  azureFunctionAppName: string;
  azureSecretRef: ResourceRef.AsObject;
  azureFunctionsList: AzureUpstreamSpec.FunctionSpec.AsObject[];
}

export const azureInitialValues: AzureValuesType = {
  azureFunctionAppName: '',
  azureSecretRef: {
    name: '',
    namespace: 'gloo-system'
  },
  azureFunctionsList: [
    {
      functionName: '',
      authLevel: AZURE_AUTH_LEVELS[0].value
    }
  ]
};

interface Props {
  parentForm: FormikProps<AzureValuesType>;
}

export const azureValidationSchema = yup.object().shape({
  azureFunctionAppName: yup.string(),
  azureSecretRefNamespace: yup.string(),
  azureSecretRefName: yup.string()
});

export const AzureUpstreamForm: React.FC<Props> = ({ parentForm }) => {
  return (
    <Formik<AzureValuesType>
      validationSchema={azureValidationSchema}
      initialValues={azureInitialValues}
      onSubmit={values => {
        parentForm.setFieldValue(
          'azureFunctionAppName',
          values.azureFunctionAppName
        );
        parentForm.setFieldValue(
          'azureFunctionsList',
          values.azureFunctionsList.slice(1)
        );
        parentForm.setFieldValue('azureSecretRef', values.azureSecretRef);
        parentForm.submitForm();
      }}>
      {({ handleSubmit }) => (
        <SoloFormTemplate formHeader='Azure Upstream Settings'>
          <InputRow>
            <Field
              name='azureFunctionAppName'
              title='Function App Name'
              placeholder='Function App Name'
              component={SoloFormInput}
            />
            <Field
              name='azureSecretRef'
              type='azure'
              component={SoloSecretRefInput}
            />
          </InputRow>
          <SoloFormTemplate formHeader='Azure Functions'>
            <FieldArray name='azureFunctionsList' render={AzureFunctions} />
          </SoloFormTemplate>
          <Footer>
            <SoloButton
              onClick={handleSubmit}
              text='Create Upstream'
              disabled={parentForm.isSubmitting}
            />
          </Footer>
        </SoloFormTemplate>
      )}
    </Formik>
  );
};

interface AzureFunctionProps extends FieldArrayRenderProps {
  form: FormikProps<AzureValuesType>;
}
const AzureFunctions: React.FC<AzureFunctionProps> = ({
  form,
  remove,
  insert,
  name
}) => {
  const cols = [
    {
      title: 'Function Name',
      dataIndex: 'functionName'
    },
    {
      title: 'Auth Level',
      dataIndex: 'lambdaFunctionName'
    },
    { title: 'Action', dataIndex: 'action' }
  ];

  const formData = form.values.azureFunctionsList.map((azureFn, index) => {
    return {
      ...azureFn,
      key: `${azureFn.functionName}`,
      action: (
        <CloseX style={{ cursor: 'pointer' }} onClick={() => remove(index)} />
      )
    };
  });

  return (
    <React.Fragment>
      <InputRow>
        <div>
          <Field
            name='azureFunctionsList[0].functionName'
            title='Function Name'
            placeholder='Function Name'
            component={SoloFormInput}
          />
        </div>
        <div>
          <Field
            name='azureFunctionsList[0].authLevel'
            title='Function Auth Level'
            defaultValue='FUNCTION'
            options={AZURE_AUTH_LEVELS}
            component={SoloFormDropdown}
          />
        </div>
        <GreenPlus
          style={{ alignSelf: 'center' }}
          onClick={() =>
            insert(0, {
              functionName: '',
              authLevel: ''
            })
          }
        />
      </InputRow>
      <InputRow>
        <Table
          style={{ width: '100%' }}
          dataSource={formData.slice(1)}
          columns={cols}
          pagination={false}
          locale={{ emptyText: '' }}
        />
      </InputRow>
    </React.Fragment>
  );
};
