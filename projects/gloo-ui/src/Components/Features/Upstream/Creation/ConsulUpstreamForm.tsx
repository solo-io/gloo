import {
  SoloFormCheckbox,
  SoloFormInput,
  SoloFormStringsList
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import * as React from 'react';
import * as yup from 'yup';

export interface ConsulVauesType {
  consulServiceName: string;
  consulServiceTagsList: string[];
  consulConnectEnabled: boolean;
  consulDataCentersList: string[];
}

export const consulInitialValues: ConsulVauesType = {
  consulServiceName: '',
  consulServiceTagsList: [''],
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
        <SoloFormTemplate formHeader='Service Tags'>
          <SoloFormStringsList
            name='consulServiceTagsList'
            title='Consul Service Tags'
            createNewPromptText='Service Tags'
          />
        </SoloFormTemplate>
      </InputRow>
      <InputRow>
        <SoloFormTemplate formHeader='Data Centers'>
          <SoloFormStringsList
            name='consulDataCentersList'
            title='Data Centers'
            createNewPromptText='Data Centers'
          />
        </SoloFormTemplate>
      </InputRow>
    </SoloFormTemplate>
  );
};
