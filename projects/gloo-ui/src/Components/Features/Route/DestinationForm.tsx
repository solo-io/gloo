import React from 'react';
import { useField, useFormikContext } from 'formik';
import {
  SoloFormDropdown,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import { UpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { HalfColumn } from './CreateRouteModal';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';
import { getFunctionList } from 'utils/helpers';

interface DestiantionFormProps {
  name: string;
  upstreamSpec: UpstreamSpec.AsObject;
}

export function DestinationForm(props: DestiantionFormProps) {
  const [field, meta] = useField(props.name);
  const { upstreamSpec } = props;
  // TODO: process upstream spec to support all types
  const functionsList = getFunctionList(props.upstreamSpec);

  return (
    <React.Fragment>
      {!!upstreamSpec && upstreamSpec.aws && (
        <React.Fragment>
          <HalfColumn>
            <SoloFormDropdown
              name={`${field.name}.aws.logicalName`}
              title='Lambda Function'
              disabled={functionsList.length === 0}
              options={functionsList}
            />
          </HalfColumn>
          <HalfColumn>
            <InputRow>
              <div>
                <SoloFormCheckbox
                  name={`${field.name}.aws.invocationStyle`}
                  title='Async'
                  disabled={functionsList.length === 0}
                />
              </div>
              <div>
                <SoloFormCheckbox
                  name={`${field.name}.aws.responseTransformation`}
                  title='Transform Response'
                  disabled={functionsList.length === 0}
                />
              </div>
            </InputRow>
          </HalfColumn>
        </React.Fragment>
      )}

      {!!upstreamSpec && upstreamSpec.kube && (
        <React.Fragment>
          <HalfColumn>
            <SoloFormDropdown
              name={`${field.name}.rest.functionName`}
              title='Function'
              disabled={functionsList.length === 0}
              options={functionsList}
            />
          </HalfColumn>
        </React.Fragment>
      )}
    </React.Fragment>
  );
}
