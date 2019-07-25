import {
  SoloFormMultipartStringCardsList,
  SoloFormInput,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import { Host } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins/static/static_pb';
import * as React from 'react';
import * as yup from 'yup';

// TODO: handle service spec
interface StaticValuesType {
  staticServiceName: string;
  staticHostList: Host.AsObject[];
  staticUseTls: boolean;
  staticServicePort: string;
}

export const staticInitialValues: StaticValuesType = {
  staticServiceName: '',
  staticHostList: [],
  staticUseTls: false,
  staticServicePort: ''
};

interface Props {}
// TODO: figure out which fields are required
export const staticValidationSchema = yup.object().shape({
  staticServicePort: yup.number(),
  staticServiceName: yup.string()
});

export const StaticUpstreamForm: React.FC<Props> = () => {
  return (
    <SoloFormTemplate formHeader='Static Upstream Settings'>
      <InputRow>
        <SoloFormInput
          name='staticServiceName'
          title='Service Name'
          placeholder='Service Name'
        />
        <SoloFormCheckbox name='staticUseTls' title='Use Tls' />
        <SoloFormInput
          name='staticServicePort'
          title='Service Port'
          placeholder='Service Port'
          type='number'
        />
      </InputRow>
      <SoloFormTemplate formHeader='Hosts'>
        <InputRow>
          <SoloFormMultipartStringCardsList
            name='staticHostList'
            createNewNamePromptText={'Address...'}
            createNewValuePromptText={'Port...'}
          />
        </InputRow>
      </SoloFormTemplate>
    </SoloFormTemplate>
  );
};
