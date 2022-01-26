import 'graphiql/graphiql.css';
import * as React from 'react';
import { GraphiQL } from 'graphiql';
import css from '@emotion/css';
import { makeExecutableSchema } from '@graphql-tools/schema';
import { GraphqlStyle } from './GraphiqlStyle';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { mapSchema, getDirective, MapperKind } from '@graphql-tools/utils';
// @ts-ignore
import { GraphQLSchema } from 'graphql';

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
  graphQLSchema: any;
}

export const GraphqlApiExplorer = (props: GraphqlApiExplorerProps) => {
  let typeDefs: any;
  if (props?.graphQLSchema?.spec?.executableSchema?.schema_definition) {
    typeDefs = props?.graphQLSchema?.spec?.executableSchema?.schema_definition;
  } else if (props?.graphQLSchema?.executableSchema?.schema_definition) {
    typeDefs = props.graphQLSchema.executableSchema.schema_definition;
  } else {
    return null;
  }
  const { mockedDirectiveTypeDefs, mockedDirectiveTransformer } =
    mockedDirective('resolve');

  let executableSchema = makeExecutableSchema({
    typeDefs: [mockedDirectiveTypeDefs],
  });

  executableSchema = mockedDirectiveTransformer(executableSchema);

  return (
    <Wrapper>
      <StyledContainer>
        {
          <GraphiQL
            css={css(GraphqlStyle)}
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
          />
        }
      </StyledContainer>
    </Wrapper>
  );
};
