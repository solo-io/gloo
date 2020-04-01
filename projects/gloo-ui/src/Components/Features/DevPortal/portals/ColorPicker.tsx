import React from 'react';
import {
  SketchPicker,
  ChromePicker,
  ColorResult,
  ColorChangeHandler
} from 'react-color';
import { css } from '@emotion/core';
import { useFormikContext, useField } from 'formik';

type ColorPickerProps = {
  initialColor?: string;
  name: string;
};

export const ColorPicker: React.FC<ColorPickerProps> = props => {
  const [displayColorPicker, setDisplayColorPicker] = React.useState(false);
  const [field, meta, helpers] = useField<string>(props.name);
  const [color, setColor] = React.useState<ColorResult>({
    hex: props.initialColor || '#B7E4CF'
  } as ColorResult);

  const handleClick = () => {
    setDisplayColorPicker(!displayColorPicker);
  };

  const handleClose = () => {
    setDisplayColorPicker(false);
  };

  const handleChange: ColorChangeHandler = color => {
    setColor(color);
    helpers.setValue(color.hex);
  };

  return (
    <div className='w-full '>
      <div
        css={css`
          background: #fff;
          border-radius: 8px;
          box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
          display: inline-block;
          cursor: pointer;
          width: 100%;
        `}
        onClick={handleClick}>
        <div className='flex items-center'>
          <div
            css={css`
              width: 36px;
              height: 35px;
              border-radius: 8px;
              border-top-right-radius: 0;
              border-bottom-right-radius: 0;
              background: ${color.hex};
            `}
          />
          <div className='ml-2'>{color.hex}</div>
        </div>
      </div>
      {displayColorPicker ? (
        <div
          css={css`
            position: absolute;
            z-index: 2;
          `}>
          <div
            css={css`
              position: fixed;
              top: 0px;
              right: 0px;
              bottom: 0px;
              left: 0px;
            `}
            onClick={handleClose}
          />
          <ChromePicker color={color.hex} {...field} onChange={handleChange} />
        </div>
      ) : null}
    </div>
  );
};
