import { graphqlConfigApi } from 'API/graphql';
import { useGetGraphqlApiYaml } from 'API/hooks';
import AreaHeader from 'Components/Common/AreaHeader';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';

const GraphqlApiConfigurationHeader: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { data: graphqlApiYaml } = useGetGraphqlApiYaml(apiRef);

  const loadYaml = async () => {
    if (!apiRef.name || !apiRef.namespace) return '';
    try {
      const yaml = await graphqlConfigApi.getGraphqlApiYaml(apiRef);
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  return (
    <div className='-mt-1 mb-5'>
      <AreaHeader
        title='Configuration'
        contentTitle={`${apiRef.namespace}--${apiRef.name}.yaml`}
        yaml={graphqlApiYaml}
        onLoadContent={loadYaml}
      />
    </div>
  );
};

export default GraphqlApiConfigurationHeader;
