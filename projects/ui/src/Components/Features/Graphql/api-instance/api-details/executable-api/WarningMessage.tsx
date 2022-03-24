import * as React from 'react';

interface Props {
  message: string;
}

export const WarningMessage: React.FC<Props> = ({ message }) => {
  return Boolean(message) ? (
    <div className='p-2 text-orange-400 border border-orange-400 mb-5 mt-5'>
      {message}
    </div>
  ) : null;
};

export default WarningMessage;
