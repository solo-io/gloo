import {
  SoloFormCheckbox,
  SoloFormDropdown
} from 'Components/Common/Form/SoloFormField';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';
import { useField } from 'formik';
import React from 'react';
import { getFunctionList } from 'utils/helpers';
import { HalfColumn } from './CreateRouteModal';
import { Upstream } from 'proto/gloo/projects/gloo/api/v1/upstream_pb';

interface DestiantionFormProps {
  name: string;
  upstream: Upstream.AsObject;
}

export function DestinationForm(props: DestiantionFormProps) {
  const [field] = useField(props.name);
  const { upstream } = props;
  // TODO: process upstream spec to support all types
  const functionsList = getFunctionList(props.upstream);

  if (functionsList.length === 0) {
    return null;
  }

  return (
    <>
      {upstream?.aws !== undefined && (
        <>
          <HalfColumn>
            <SoloFormDropdown
              name={`${field.name}.aws.logicalName`}
              title='Lambda Function'
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
        </>
      )}

      {upstream?.kube !== undefined && (
        <>
          <HalfColumn>
            <SoloFormDropdown
              name={`${field.name}.rest.functionName`}
              title='Function'
              disabled={functionsList.length === 0}
              options={functionsList}
            />
          </HalfColumn>
        </>
      )}
    </>
  );
}
