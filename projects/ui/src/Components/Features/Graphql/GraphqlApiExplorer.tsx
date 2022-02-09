import * as React from 'react';
import { GraphiQL } from 'graphiql';
import { makeExecutableSchema } from '@graphql-tools/schema';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { mapSchema, getDirective, MapperKind } from '@graphql-tools/utils';
// @ts-ignore
import { GraphQLSchema } from 'graphql';
import { useListVirtualServices } from 'API/hooks';
import { useParams } from 'react-router';
import { VirtualService } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import {
  OverviewSmallBoxSummary,
  StatusHealth,
  WarningCircle,
} from '../Overview/OverviewBoxSummary';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';

function mockedDirective(directiveName: string) {
  return {
    mockedDirectiveTypeDefs: `directive @${directiveName}(name: String) on FIELD_DEFINITION | ENUM_VALUE | INPUT_FIELD_DEFINITION | INPUT_OBJECT | OBJECT | SCALAR | ARGUMENT_DEFINITION `,
    mockedDirectiveTransformer: (schema: GraphQLSchema) =>
      mapSchema(schema, {
        [MapperKind.OBJECT_FIELD]: fieldConfig => {
          const mockedDirective = getDirective(
            schema,
            fieldConfig,
            directiveName
          )?.[0];
          if (mockedDirective) {
            fieldConfig.deprecationReason = mockedDirective['name'];
            return fieldConfig;
          }
        },
        [MapperKind.ENUM_VALUE]: enumValueConfig => {
          const mockedDirective = getDirective(
            schema,
            enumValueConfig,
            directiveName
          )?.[0];
          if (mockedDirective) {
            enumValueConfig.deprecationReason = mockedDirective['name'];
            return enumValueConfig;
          }
        },
      }),
  };
}
const Wrapper = styled.div`
  background: white;
`;

const Header = styled.h1`
  background: ${colors.marchGrey};
  padding: 20px;
  margin-bottom: 0;
  border-radius: 10px 10px 0 0;
`;

const StyledContainer = styled.div`
  height: 70vh;
`;

interface GraphqlApiExplorerProps {
  graphQLSchema?: any;
}

export const GraphqlApiExplorer = (props: GraphqlApiExplorerProps) => {
  const { name, namespace } = useParams();
  const graphiqlRef = React.useRef<null | GraphiQL>(null);

  const { mockedDirectiveTypeDefs, mockedDirectiveTransformer } =
    mockedDirective('resolve');
  const { data: virtualServices, error: virtualServicesError } =
    useListVirtualServices();
  const [correspondingVirtualServices, setCorrespondingVirtualServices] =
    React.useState<VirtualService.AsObject[]>([]);

  React.useEffect(() => {
    let correspondingVs = virtualServices?.filter(vs =>
      vs.spec?.virtualHost?.routesList.some(
        route =>
          route?.graphqlSchemaRef?.name === name &&
          route?.graphqlSchemaRef?.namespace === namespace
      )
    );

    if (!!correspondingVs) {
      setCorrespondingVirtualServices(correspondingVs);
    }
  }, []);

  let executableSchema = makeExecutableSchema({
    typeDefs: [mockedDirectiveTypeDefs],
  });

  executableSchema = mockedDirectiveTransformer(executableSchema);

  const handlePrettifyQuery = () => {
    graphiqlRef?.current?.handlePrettifyQuery();
  };

  // TODO:  We can hide and show elements based on what we get back.

  return !!correspondingVirtualServices?.length ? (
    <Wrapper>
      <StyledContainer>
        <GraphiQL
          ref={graphiqlRef}
          schema={executableSchema}
          fetcher={async graphQLParams => {
            try {
              const data = await fetch('http://localhost:8080/graphql', {
                method: 'POST',
                headers: {
                  Accept: 'application/json',
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify(graphQLParams),
                credentials: 'same-origin',
              });
              return data.json().catch(() => data.text());
            } catch (error) {
              console.log('error', error);
            }
          }}
        >
          <GraphiQL.Toolbar>
            <GraphiQL.Button
              onClick={handlePrettifyQuery}
              label='Prettify'
              title='Prettify Query (Shift-Ctrl-P)'
            />
          </GraphiQL.Toolbar>
        </GraphiQL>
      </StyledContainer>
    </Wrapper>
  ) : (
    <div className='overflow-hidden bg-white rounded-lg shadow'>
      <div className='px-4 py-5 sm:p-6'>
        <StatusHealth isWarning>
          <div>
            <WarningCircle>
              <WarningExclamation />
            </WarningCircle>
          </div>
          <div>
            <>
              <div className='text-xl '>Explorer unavailable</div>
              <div className='text-lg '>
                There is no Virtual Service that exposes this GraphQL endpoint
              </div>
            </>
          </div>
        </StatusHealth>
      </div>
    </div>
  );
};
