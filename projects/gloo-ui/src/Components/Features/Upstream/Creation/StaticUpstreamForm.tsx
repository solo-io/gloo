import {
  SoloFormMultipartStringCardsList,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import {
  InputRow,
  SoloFormTemplate
} from 'Components/Common/Form/SoloFormTemplate';
import * as React from 'react';
import * as yup from 'yup';
import { useFormikContext } from 'formik';

// TODO: handle service spec
export interface StaticValuesType {
  staticServiceName: string;
  staticHostList: { name: string; value: string }[];
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
const staticValidationSchema = yup.object().shape({
  staticServicePort: yup.number(),
  staticServiceName: yup.string()
});

export const StaticUpstreamForm: React.FC<Props> = () => {
  const form = useFormikContext<StaticValuesType>();

  return (
    <SoloFormTemplate formHeader='Static Upstream Settings'>
      <InputRow>
        <SoloFormMultipartStringCardsList
          name='staticHostList'
          title='Hosts'
          values={form.values.staticHostList}
          createNewNamePromptText={'address...'}
          createNewValuePromptText={'port...'}
        />

        <SoloFormCheckbox name='staticUseTls' title='Use Tls' />
      </InputRow>
      <InputRow />
    </SoloFormTemplate>
  );
};
