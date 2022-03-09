import styled from '@emotion/styled/macro';
import { useGetGraphqlSchemaDetails } from 'API/hooks';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import React from 'react';
import { useParams } from 'react-router';
import { useVirtual } from 'react-virtual';
import { colors } from 'Styles/colors';
import tw from 'twin.macro';
import gql from 'graphql-tag';
import {
  EnumValueDefinitionNode,
  FieldDefinitionNode,
  NamedTypeNode,
  //@ts-ignore
} from 'graphql';
import { GraphqlIconHolder } from './../GraphqlTable';
import { OperationDescription } from './../ResolversTable';

export const EnumResolver: React.FC<{
    resolverType: string;
    fields: EnumValueDefinitionNode[];
  }> = props => {
    const { resolverType, fields } = props;

    const listRef = React.useRef<HTMLDivElement>(null);

    const rowVirtualizer = useVirtual({
      size: fields?.length ?? 0,
      parentRef: listRef,
      estimateSize: React.useCallback(() => 90, []),
      overscan: 1,
    });

    return (
      <div key={resolverType}>
        <div className='relative flex flex-col w-full bg-gray-200 border h-28'>
          <div className='flex items-center justify-between gap-5 pt-4 my-2 ml-4 h-14 '>
            <div className='flex items-center mr-3'>
              <GraphqlIconHolder>
                <GraphQLIcon className='w-4 h-4 fill-current' />
              </GraphqlIconHolder>
              <span className='flex items-center font-medium text-gray-900 whitespace-nowrap'>
                {resolverType}
              </span>
            </div>
          </div>
          <div className='flex items-center justify-between w-full px-6 py-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
            <div
              className='relative flex-wrap justify-between w-full h-full text-sm '
              style={{
                display: 'grid',
                flexWrap: 'wrap',
                gridTemplateColumns: '1fr 1fr  1fr',
                gridTemplateRows: '1fr',
                gridAutoRows: 'min-content',
                columnGap: '15px',
              }}>
              <span className='flex items-center justify-start ml-6 font-medium text-gray-900 '>
                Value
              </span>
              <span className='flex items-center justify-start ml-8 font-medium text-gray-900 '>
                Description
              </span>


            </div>
          </div>

        </div>

          <div
            ref={listRef}
            style={{
              height: `${
                fields?.length * 90 < 400 ? fields!.length * 90 : 400
              }px`,
              width: `100%`,
              overflow: 'auto',
            }}>
            <div
              style={{
                height: `${rowVirtualizer.totalSize}px`,
                width: '100%',
                position: 'relative',
              }}>
              {rowVirtualizer.virtualItems.map(virtualRow => {
                const op = fields[virtualRow.index] as EnumValueDefinitionNode;
                return (
                  <div
                    key={`${resolverType}-${op.name?.value}`}
                    className={`flex h-20 p-2 pl-0 border `}
                    style={{
                      position: 'absolute',
                      top: 0,
                      left: 0,
                      width: '100%',
                      height: `${virtualRow.size}px`,
                      transform: `translateY(${virtualRow.start}px)`,
                    }}>
                    <div className='flex items-center px-4 text-sm font-medium text-gray-900 whitespace-nowrap'>
                      <CodeIcon className='w-4 h-4 ml-2 mr-3 fill-current text-blue-600gloo' />
                    </div>
                    <div className='relative flex items-center w-full text-sm text-gray-500 whitespace-nowrap'>
                      <div
                        className='relative flex-wrap justify-between w-full h-full text-sm '
                        style={{
                          display: 'grid',
                          flexWrap: 'wrap',
                          gridTemplateColumns:
                            '1fr 1fr  minmax(120px, 250px) 105px',
                          gridTemplateRows: op.description?.value
                            ? ' 1fr min-content'
                            : '1fr',
                          gridAutoRows: 'min-content',
                          columnGap: '5px',
                          rowGap: '5px',
                        }}>
                        <span className='flex items-center font-medium text-gray-900 '>
                          {op.name.value}
                        </span>
                        <span className='flex items-center text-sm text-gray-700 '>
                          {op.description?.value}
                        </span>



                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
      </div>
    );
  };
