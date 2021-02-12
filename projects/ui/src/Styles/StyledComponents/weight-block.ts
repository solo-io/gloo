import styled from '@emotion/styled/macro';
import { colors } from 'Styles/colors';

type WeightPercentageBlockProps = {
  percentage?: number;
  width?: string;
};

export const WeightPercentageBlock = styled.div<WeightPercentageBlockProps>`
  color: white;
  width: ${(props: WeightPercentageBlockProps) => props.width ?? '105px'};
  border-radius: 4px;
  padding: 2px;
  text-align: center;

  background: ${(props: WeightPercentageBlockProps) =>
    props.percentage === undefined
      ? colors.juneGrey
      : props.percentage > 80
      ? colors.planeOfWaterBlue
      : props.percentage > 50
      ? colors.neptuneBlue
      : props.percentage > 30
      ? colors.oceanBlue
      : props.percentage > 10
      ? colors.seaBlue
      : colors.pondBlue};
`;
