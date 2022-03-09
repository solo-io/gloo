/* 
Based on this example:
https://tailwindui.com/components/application-ui/navigation/navbars
*/
import { Disclosure, Menu, Transition } from '@headlessui/react';
import React from 'react';
import { CloseOutlined, MenuOutlined } from '@ant-design/icons';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { Link } from 'react-router-dom';
import { SoloNavbarLink } from './SoloNavbarLink';

const StyledDisclosure = styled(Disclosure)`
  background-color: ${colors.seaBlue};
  z-index: 1;
`;
const StyledMobileNav = styled.div`
  background-color: ${colors.seaBlue};
`;
const StyledVerticalSeparator = styled.div`
  margin-left: 50px;
  height: 70%;
  border-right: 2px solid ${colors.lakeBlue};
  margin-right: 50px;
`;

export type ISoloNavLink = { name: string; href: string; exact?: boolean };
interface ISoloNavbar {
  BrandComponent?: React.FunctionComponent;
  navLinks: ISoloNavLink[];
  SettingsComponent?: React.FunctionComponent;
}
const SoloNavbar: React.FC<ISoloNavbar> = ({
  BrandComponent,
  navLinks,
  SettingsComponent,
}) => {
  return (
    <StyledDisclosure as='nav'>
      {({ open }) => (
        <>
          <div className='max-w-7xl mx-auto px-2 lg:px-6 xl:px-8'>
            <div className='relative flex items-center justify-between h-[55px]'>
              <div className='absolute inset-y-0 left-0 flex items-center lg:hidden'>
                {/* Mobile menu button*/}
                <Disclosure.Button className='inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white'>
                  <span className='sr-only'>Open main menu</span>
                  {open ? (
                    <CloseOutlined
                      className='block h-6 w-6 text-[1.5rem]'
                      aria-hidden='true'
                    />
                  ) : (
                    <MenuOutlined
                      className='block h-6 w-6 text-[1.5rem]'
                      aria-hidden='true'
                    />
                  )}
                </Disclosure.Button>
              </div>
              <div className='flex-1 flex items-center justify-center lg:items-stretch lg:justify-start h-[100%]'>
                <div className='flex-shrink-0 flex items-center'>
                  {BrandComponent && (
                    <>
                      <Link to='/'>
                        <BrandComponent />
                      </Link>
                      <StyledVerticalSeparator />
                    </>
                  )}
                </div>
                <div className='hidden lg:block'>
                  <div className='flex h-[100%] items-end'>
                    {navLinks.map(item => (
                      <div key={item.name} className='mr-[50px]'>
                        <SoloNavbarLink to={item.href} exact={item.exact}>
                          {item.name}
                        </SoloNavbarLink>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <div className='absolute inset-y-0 right-0 flex items-center pr-2 lg:static lg:inset-auto lg:ml-6 lg:pr-0'>
                {SettingsComponent && <SettingsComponent />}
              </div>
            </div>
          </div>

          <Disclosure.Panel className='xl:hidden'>
            <StyledMobileNav className='px-2 pt-2 pb-3 space-y-1'>
              {navLinks.map(item => (
                <div key={item.name} className='block px-3 py-2'>
                  <Disclosure.Button>
                    <SoloNavbarLink to={item.href} exact={item.exact}>
                      {item.name}
                    </SoloNavbarLink>
                  </Disclosure.Button>
                </div>
              ))}
            </StyledMobileNav>
          </Disclosure.Panel>
        </>
      )}
    </StyledDisclosure>
  );
};

export default SoloNavbar;
