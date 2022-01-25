import * as React from 'react';
import styled from '@emotion/styled/macro';
import { SoloModal } from 'Components/Common/SoloModal';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloInput } from 'Components/Common/SoloInput';
import { colors } from 'Styles/colors';
// import { SoloFi}

export interface NewApiModalProps {
    showNewModal: boolean;
    toggleNewModal: () => any;
}

const InputWrapper = styled.div`
    padding: 10px 0;
`;

const Footer = styled.footer`
    display: flex;
    flex-direction: row-reverse;
`;

const Button = styled.button`
    background-color: ${colors.seaBlue};
    color: white;
    padding: 15px;
    border: none;
    &:hover {
        cursor: pointer;
    }
`;

export const NewApiModal = (props: NewApiModalProps) => {
    const { showNewModal, toggleNewModal } = props;
    const [name, setName] = React.useState<string>('');
    const [_schemaFile, setSchemaFile] = React.useState<File>();

    const changeName = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value);
    }

    // Check .graphql files as well.
    const changeSchema = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setSchemaFile(file);
        };
    }

    return (
        <SoloModal visible={showNewModal} width={1000} onClose={toggleNewModal}>
            <SectionCard cardName="Create new GraphQL API">
                <InputWrapper>
                    <SoloInput title='Name' onChange={changeName} value={name} />
                </InputWrapper>
                <InputWrapper>
                    {/* @ts-expect-error Setting the value here will cause an error. */}
                    <SoloInput title='Schema Definition' file onChange={changeSchema} />
                </InputWrapper>
                <Footer>
                    <Button type="button" onClick={toggleNewModal}>
                        Create API
                    </Button>
                </Footer>
            </SectionCard>
        </SoloModal>
    )
}
