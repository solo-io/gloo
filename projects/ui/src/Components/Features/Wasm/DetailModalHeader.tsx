import React from 'react';
import styled from '@emotion/styled';
import { CardHeader } from 'Components/Common/Card';
import { colors } from 'Styles/colors';

export const PolicyHeaderTitle = styled.div`
  display: flex;
  align-items: center;
  font-weight: 500;
  font-size: 18px;
  margin-left: 23px;
`;

export const PolicyHeaderData = styled.div`
  font-size: 24px;
  line-height: 34px;
  border: 1px solid ${colors.aprilGrey};
  border-radius: 8px;
  font-weight: 600;
  margin-left: 10px;
  padding: 0 8px;
`;

const PolicyDetailsHeader = styled(CardHeader)`
  background: white;
  justify-content: space-between;
  padding-right: 40px;
  line-height: normal;
  padding-top: 20px;

  > div {
    display: flex;
    align-items: center;
  }
`;

const HeaderImageHolder = styled.div`
  margin-right: 15px;
  height: 33px;
  width: 33px;
  border-radius: 100%;
  background: ${colors.seaBlue};
  display: flex;
  justify-content: center;
  align-items: center;

  img,
  svg {
    width: 30px;
    max-height: 30px;

    * {
      fill: white;
    }
  }
`;

const HeaderTitleSection = styled.div`
  max-width: calc(100% - 300px);
  min-height: 58px;
`;
const HeaderTitleName = styled.div`
  width: 100%;
  font-size: 22px;
  color: ${colors.novemberGrey};
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
`;
const SecondaryInformation = styled.div`
  display: flex;
  align-items: center;

  font-size: 14px;
  line-height: 22px;
  height: 22px;
  padding: 0 12px;
  color: ${colors.novemberGrey};
  background: white;
  margin-left: 13px;
  border-radius: 16px;
`;

const MainSection = styled.div`
  padding: 0 13px;
`;

const Divider = styled.div`
  height: 1px;
  background: ${colors.februaryGrey};
  margin: 18px 0;
`;

const PolicyDetailsDescription = styled.div`
  border-radius: 8px;
  background: ${colors.januaryGrey};
  padding: 15px 13px;
  margin-bottom: 18px;
`;

type HeaderProps = {
  headerName: string;
  logoIcon?: React.ReactNode;
  secondaryComponent?: React.ReactNode;
  description?: React.ReactNode;
};
export const DetailModalHeader = ({
  headerName,
  logoIcon,
  secondaryComponent,
  description,
}: HeaderProps) => {
  return (
    <>
      <PolicyDetailsHeader>
        <div>
          {logoIcon && <HeaderImageHolder>{logoIcon}</HeaderImageHolder>}
          <div>
            <HeaderTitleName>{headerName}</HeaderTitleName>
          </div>
        </div>
        {!!secondaryComponent && (
          <SecondaryInformation>{secondaryComponent}</SecondaryInformation>
        )}
      </PolicyDetailsHeader>
      <MainSection>
        <Divider />
        <PolicyDetailsDescription>{description}</PolicyDetailsDescription>
      </MainSection>
    </>
  );
};
