import styled from '@emotion/styled/macro';

interface IconProps {
  width?: number;
  applyColor?: {
    strokeNotFill?: boolean;
    color: string;
  };
}

export const IconHolder = styled.div<IconProps>`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    ${(props: IconProps) => `width: ${props.width || 25}px;`}

    ${(props: IconProps) =>
      props.applyColor
        ? `
            * {
              ${props.applyColor.strokeNotFill ? 'stroke' : 'fill'}: ${
            props.applyColor.color
          };
            }`
        : ''}
  }
`;

type CircleIconHolderProps = {
  backgroundColor?: string;
  iconColor?: {
    strokeNotFill?: boolean;
    color: string;
  };
  iconSize?: string;
};

export const CircleIconHolder = styled.div<CircleIconHolderProps>`
  height: 33px;
  width: 33px;
  border-radius: 100%;
  display: flex;
  justify-content: center;
  align-items: center;

  background: ${(props: CircleIconHolderProps) =>
    props.backgroundColor ?? 'white'};

  svg {
    width: ${(props: CircleIconHolderProps) => props.iconSize ?? '30px'};
    height: ${(props: CircleIconHolderProps) => props.iconSize ?? '30px'};
    ${(props: CircleIconHolderProps) =>
      props.iconColor
        ? `
            * {
              ${props.iconColor.strokeNotFill ? 'stroke' : 'fill'}: ${
            props.iconColor.color
          };
            }`
        : ''}
  }
`;
