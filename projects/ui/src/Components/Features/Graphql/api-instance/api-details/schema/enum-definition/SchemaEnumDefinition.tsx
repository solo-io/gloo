import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { EnumValueDefinitionNode } from 'graphql';
import React from 'react';
import { useVirtual } from 'react-virtual';
import { colors } from 'Styles/colors';

export const SchemaEnumDefinition: React.FC<{
  resolverType: string;
  values: readonly EnumValueDefinitionNode[];
}> = ({ resolverType, values }) => {
  const listRef = React.useRef<HTMLDivElement>(null);

  const rowVirtualizer = useVirtual({
    size: values?.length ?? 0,
    parentRef: listRef,
    estimateSize: React.useCallback(() => 90, []),
    overscan: 1,
  });

  return (
    <div data-testid='enum-resolver' key={resolverType}>
      <div
        className='relative flex flex-col w-full py-3 border'
        style={{
          backgroundColor: colors.lightJanuaryGrey,
          display: 'grid',
          flexWrap: 'wrap',
          gridTemplateColumns: '1fr 1fr  1fr',
          gridTemplateRows: '1fr',
          gridAutoRows: 'min-content',
          columnGap: '15px',
        }}>
        <span className='flex items-center justify-start ml-8 font-medium text-gray-900 '>
          Value
        </span>
        <span className='flex items-center justify-center ml-8 font-medium text-gray-900 '>
          Description
        </span>
      </div>

      <div
        ref={listRef}
        style={{
          height: `${values?.length * 90 < 400 ? values!.length * 90 : 400}px`,
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
            const op = values[virtualRow.index] as EnumValueDefinitionNode;
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
