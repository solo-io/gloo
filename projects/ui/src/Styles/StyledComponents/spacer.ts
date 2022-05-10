import styled from '@emotion/styled';

const spacerScale = '.25rem';

/**
 * Use a number (e.g. `mx={3}`) for scaled default spacing,
 * or a string for exact values (e.g. `mx='5px'`).
 * - Note: For overrides, `mr`, `mb` and other
 * directional styles override `mx`, `my`.
 */
export const Spacer = styled.div<{
  mx?: string | number;
  my?: string | number;
  mr?: string | number;
  ml?: string | number;
  mt?: string | number;
  mb?: string | number;
  margin?: string | number;
  px?: string | number;
  py?: string | number;
  pr?: string | number;
  pl?: string | number;
  pt?: string | number;
  pb?: string | number;
  padding?: string | number;
}>(props => {
  //
  // Scale any numeric values with `spacerScale`.
  Object.keys(props).forEach(k => {
    const value = (props as any)[k] as string | number;
    if (typeof value === 'number')
      (props as any)[k] = `calc(${value} * ${spacerScale})`;
  });
  //
  // Check for overrides, return styles.
  const { mx, my, mr, ml, mt, mb, margin, px, py, pl, pr, pt, pb, padding } =
    props;
  const marginRight = mr ?? mx;
  const marginLeft = ml ?? mx;
  const marginTop = mt ?? my;
  const marginBottom = mb ?? my;
  const paddingRight = pr ?? px;
  const paddingLeft = pl ?? px;
  const paddingTop = pt ?? py;
  const paddingBottom = pb ?? py;
  return `
  ${margin ? 'margin:' + margin + ';' : ''}
  ${marginRight ? 'margin-right:' + marginRight + ';' : ''}
  ${marginLeft ? 'margin-left:' + marginLeft + ';' : ''}
  ${marginTop ? 'margin-top:' + marginTop + ';' : ''}
  ${marginBottom ? 'margin-bottom:' + marginBottom + ';' : ''}
  ${padding ? 'padding:' + padding + ';' : ''}
  ${paddingRight ? 'padding-right:' + paddingRight + ';' : ''}
  ${paddingLeft ? 'padding-left:' + paddingLeft + ';' : ''}
  ${paddingTop ? 'padding-top:' + paddingTop + ';' : ''}
  ${paddingBottom ? 'padding-bottom:' + paddingBottom + ';' : ''}
`;
});
