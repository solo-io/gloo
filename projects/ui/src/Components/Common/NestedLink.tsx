import React from 'react';
import { useLocation, Link, LinkProps } from 'react-router-dom';
import { colors } from 'Styles/colors';

const NavLinkStyles = {
  display: 'inline-block',
  color: 'white',
  textDecoration: 'none',
  fontSize: '18px',
  marginRight: '50px',
  fontWeight: 300,
  borderBottom: `8px solid transparent`
};
const activeStyle = {
  borderBottom: `8px solid ${colors.pondBlue}`,
  cursor: 'default',
  fontWeight: 400
};

export const NestedLink: React.FC<
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
