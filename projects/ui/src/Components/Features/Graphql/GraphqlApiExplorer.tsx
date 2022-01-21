import 'graphiql/graphiql.css';
import * as React from 'react';
import { GraphiQL } from 'graphiql';
import css from '@emotion/css';
import { createGraphiQLFetcher } from '@graphiql/toolkit';
import { makeExecutableSchema } from '@graphql-tools/schema';
import { GraphqlStyle } from './GraphiqlStyle';
import styled from '@emotion/styled';
import { useApiProvider } from './state/ApiProvider.state';
import { colors } from 'Styles/colors';

const Wrapper = styled.div`
    margin: 13px;
    background: white;
    border-radius: 10px;
    box-shadow: 0px 4px 9px ${colors.boxShadow};
`;

const Header = styled.h1`
    background: ${colors.marchGrey};
    padding: 20px;
    margin-bottom: 0;
    border-radius: 10px 10px 0 0;
`

const StyledContainer = styled.div`
    height: 100vh;

`;

interface GraphqlApiExplorerProps {
    graphQLSchema: any;
}

export const GraphqlApiExplorer = (props: GraphqlApiExplorerProps) => {
    const { state } = useApiProvider();
    let typeDefs: any;
    if (props?.graphQLSchema?.spec?.executableSchema?.schema_definition) {
        typeDefs = props?.graphQLSchema?.spec?.executableSchema?.schema_definition;
    } else if (props?.graphQLSchema?.executableSchema?.schema_definition) {
        typeDefs = props.graphQLSchema.executableSchema.schema_definition;
    } else {
        return null;
    }

    const executableSchema = makeExecutableSchema({
        typeDefs,
    });

    const fetcher = createGraphiQLFetcher({
        url: state.url as string,
    });
    return (
        <Wrapper>
            <StyledContainer>
                <GraphiQL css={css(GraphqlStyle)} schema={executableSchema} fetcher={fetcher} />
            </StyledContainer>
        </Wrapper>
    )
}
