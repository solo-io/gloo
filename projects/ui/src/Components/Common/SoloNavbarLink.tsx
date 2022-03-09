import React, { CSSProperties } from 'react';
import { useLocation, Link, LinkProps } from 'react-router-dom';
import { colors } from 'Styles/colors';

const NavLinkStyles = {
  display: 'inline-block',
  color: 'white',
  textDecoration: 'none',
  fontSize: '18px',
  whiteSpace: 'nowrap',
  fontWeight: 300,
  borderBottom: `8px solid transparent`,
  paddingBottom: '.5rem',
} as CSSProperties;
const activeStyle = {
  borderBottom: `8px solid ${colors.pondBlue}`,
  fontWeight: 400,
} as CSSProperties;

export const SoloNavbarLink: React.FC<
  LinkProps & { exact?: boolean; isActive?: (pathname: string) => boolean }
> = ({ to, children, exact, isActive, ...rest }) => {
  const routerLocation = useLocation();
  let linkStyles = { ...NavLinkStyles };

  if (exact) {
    if (routerLocation.pathname === to.toString()) {
      linkStyles = { ...linkStyles, ...activeStyle };
    }
  } else if (isActive) {
    if (isActive(routerLocation.pathname)) {
      linkStyles = { ...linkStyles, ...activeStyle };
    }
  } else if (routerLocation.pathname.includes(to.toString())) {
    linkStyles = { ...linkStyles, ...activeStyle };
  }

  return (
    <Link {...rest} to={to} style={linkStyles}>
      {children}
    </Link>
  );
};
