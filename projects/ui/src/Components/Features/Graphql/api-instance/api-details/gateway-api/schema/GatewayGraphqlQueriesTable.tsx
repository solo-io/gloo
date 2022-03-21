import { SoloInput } from 'Components/Common/SoloInput';
import { SoloTable } from 'Components/Common/SoloTable';
import React, { useMemo, useState } from 'react';

const GWGraphqlQueriesTable = () => {
  const tableData: any = [];

  const columns = useMemo(() => {
    return [
      {
        title: 'Query Name',
        dataIndex: 'name',
        width: 200,
        // render: RenderSimpleLink,
      },
      {
        title: 'Namespace',
        dataIndex: 'namespace',
      },
      {
        title: 'Resolvers',
        dataIndex: 'resolvers',
      },
    ];
  }, []);

  const [searchText, setSearchText] = useState('');
  return (
    <div>
      <div className='w-[400px] mb-5'>
        <SoloInput
          value={searchText}
          placeholder='Search by name...'
          onChange={s => setSearchText(s.target.value)}
        />
      </div>
      <SoloTable columns={columns} dataSource={tableData} />
      {/* <ConfirmationModal
        visible={isDeleting}
        confirmPrompt='delete this API'
        confirmButtonText='Delete'
        goForIt={deleteFn}
        cancel={cancelDelete}
        isNegative
      />
      <ErrorModal
        {...errorDeleteModalProps}
        cancel={closeErrorModal}
        visible={errorModalIsOpen}
      /> */}
    </div>
  );
};

export default GWGraphqlQueriesTable;
