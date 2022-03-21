import React from 'react';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';

const SoloAddButton: React.FC<
  React.ButtonHTMLAttributes<HTMLButtonElement>
> = props => (
  <button className='inline-block' {...props}>
    <span className='flex items-center text-green-400 cursor-pointer hover:text-green-300'>
      <GreenPlus className='w-6 mr-1 fill-current' />
      <span className='text-gray-700'>&nbsp;{props.children}</span>
    </span>
  </button>
);

export default SoloAddButton;
