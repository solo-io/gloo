import {
  EnumTypeDefinitionNode,
  Kind,
  ObjectTypeDefinitionNode,
} from 'graphql';

export const mockEnumDefinitions = [
  {
    kind: Kind.ENUM_TYPE_DEFINITION,
    name: { kind: Kind.NAME, value: 'test enum 1' },
    values: [
      {
        kind: Kind.ENUM_VALUE_DEFINITION,
        description: {
          kind: Kind.STRING,
          value: 'test description a',
        },
        name: { kind: Kind.NAME, value: 'test value a' },
      },
      {
        kind: Kind.ENUM_VALUE_DEFINITION,
        description: {
          kind: Kind.STRING,
          value: 'test description b',
        },
        name: { kind: Kind.NAME, value: 'test value b' },
      },
    ],
  },
  {
    kind: Kind.ENUM_TYPE_DEFINITION,
    name: { kind: Kind.NAME, value: 'test enum 2' },
    values: [
      {
        kind: Kind.ENUM_VALUE_DEFINITION,
        description: {
          kind: Kind.STRING,
          value: 'test description 2',
        },
        name: { kind: Kind.NAME, value: 'test value 2' },
      },
    ],
  },
] as (EnumTypeDefinitionNode | ObjectTypeDefinitionNode)[];
