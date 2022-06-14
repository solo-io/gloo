import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GrpcProtoCheck } from '../../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/grpcProtoCheck/GrpcProtoCheck';
import { Formik } from 'formik';

export default {
  title:
    'Graphql / api-instance / api-details / schema / resolver-wizard / grpcProtoCheck / GrpcProtoCheck',
  component: GrpcProtoCheck,
} as ComponentMeta<typeof GrpcProtoCheck>;

const Template: ComponentStory<typeof GrpcProtoCheck> = args => (
  <Formik
    onSubmit={() => {}}
    initialValues={{
      resolverType: {},
      upstream: '',
      resolverConfig: '',
      listOfResolvers: [],
      protoFile: '',
    }}>
    <GrpcProtoCheck {...args} />
  </Formik>
);

export const Primary = Template.bind({});

// @ts-ignore
Primary.args = {
  warningMessage: '',
} as Partial<typeof GrpcProtoCheck>;
