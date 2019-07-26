import React from 'react';
import { useField } from 'formik';
import {
  SoloFormDropdown,
  SoloFormCheckbox
} from 'Components/Common/Form/SoloFormField';
import { UpstreamSpec } from 'proto/github.com/solo-io/gloo/projects/gloo/api/v1/plugins_pb';
import { HalfColumn } from './CreateRouteModal';
import { InputRow } from 'Components/Common/Form/SoloFormTemplate';

interface DestiantionFormProps {
  name: string;
  upstreamSpec: UpstreamSpec.AsObject;
}

export function DestinationForm(props: DestiantionFormProps) {
  const [field, meta] = useField(props.name);

  // TODO: process upstream spec to support all types
  const [upstreamSpec, setUpstreamSpec] = React.useState<
    UpstreamSpec.AsObject
  >();
  const [functionsList, setFunctionsList] = React.useState<
    { key: string; value: string }[]
  >([]);

  React.useEffect(() => {
    if (props.upstreamSpec) {
      setUpstreamSpec(props.upstreamSpec);
      if (props.upstreamSpec.aws) {
        let newList = props.upstreamSpec.aws.lambdaFunctionsList.map(lambda => {
          return {
            key: lambda.logicalName,
            value: lambda.logicalName
          };
        });
        setFunctionsList(newList);
      }
    }
  }, [props.upstreamSpec]);

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
          <InputRow>
            <HalfColumn>
              <SoloFormCheckbox
                name={`${field.name}.aws.invocationStyle`}
                title='Async'
                disabled={functionsList.length === 0}
              />
            </HalfColumn>
            <HalfColumn>
              <SoloFormCheckbox
                name={`${field.name}.aws.responseTransformation`}
                title='Transform Response'
                disabled={functionsList.length === 0}
              />
            </HalfColumn>
          </InputRow>
        </React.Fragment>
      )}
    </React.Fragment>
  );
}
