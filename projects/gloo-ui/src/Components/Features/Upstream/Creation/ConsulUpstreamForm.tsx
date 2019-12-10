import {
  SoloFormCheckbox,
  SoloFormInput,
  SoloFormStringsList,
  ErrorText
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import * as React from 'react';
import * as yup from 'yup';
import { Select } from 'antd';
import { useFormikContext } from 'formik';
import { Label } from 'Components/Common/SoloInput';
const { Option } = Select;
export interface ConsulVauesType {
  consulServiceName: string;
  consulServiceTagsList: string[];
  consulConnectEnabled: boolean;
  consulDataCentersList: string[];
}

export const consulInitialValues: ConsulVauesType = {
  consulServiceName: '',
  consulServiceTagsList: [] as string[],
  consulConnectEnabled: false,
  consulDataCentersList: ['']
};

interface Props {}

export const consulValidationSchema = yup.object().shape({
  consulServiceName: yup.string(),
  consulServiceTagsList: yup.array().of(yup.string()),
  consulConnectEnabled: yup.boolean(),
  consulDataCentersList: yup.array().of(yup.string())
});

export const SoloFormTagSelect = ({ ...props }) => {
  const { name, label, placeholder } = props;
  const form = useFormikContext<ConsulVauesType>();
  const field = form.getFieldProps(name);
  const meta = form.getFieldMeta(name);

  function handleChange(newVal: any) {}

  return (
    <>
      {label && <Label>{label}</Label>}
      <Select
        {...field}
        {...props}
        mode='tags'
        size={'default'}
        placeholder={placeholder}
        onChange={handleChange}
        style={{ width: '100%' }}>
        {field.value.map((opt: any) => (
          <Option key={opt}>{opt}</Option>
        ))}
      </Select>
      <ErrorText errorExists={!!meta.error && meta.touched}>
        {meta.error}
      </ErrorText>
    </>
  );
};

export const ConsulUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Consul Upstream Settings'>
      <InputRow>
        <div>
          <SoloFormInput
            name='consulServiceName'
            title='Service Name'
            placeholder='Service Name'
          />
        </div>
        <div>
          <SoloFormCheckbox
            name='consulConnectEnabled'
            title='Enable Consul Connect'
          />
        </div>
      </InputRow>
      <InputRow>
        <div>
          <SoloFormStringsList
            data-testid='consul-service-tags'
            name='consulServiceTagsList'
            label='Consul Service Tags'
            createNewPromptText='Service Tags'
          />
        </div>
        <div>
          <SoloFormStringsList
            label='Data Centers'
            name='consulDataCentersList'
            title='Data Centers'
            createNewPromptText='Data Centers'
          />
        </div>
      </InputRow>
    </SoloFormTemplate>
  );
};
