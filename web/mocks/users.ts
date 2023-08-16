import { generateName, randomInt } from './utils';

/**
 * Generates a mock users. For the purposes of the team id,
 * five possible teams are assumed unless otherwise specified.
 * @param teams The number of available teams.
 */
const createUser = ( teams?: number ): IUser => {
  const lastName = generateName( 10 );

  return {
    firstName: generateName( 10 ),
    lastName,
    email: `${lastName.toLowerCase()}@test.com`,
    team: ( randomInt( teams || 5 ) + 1 ).toString(),
    active: randomInt( 2 ) === 0,
  };
};

/**
 * Generates a list of users of a specified length.
 * @param count The number of teams to create.
 * @param teams The number of available teams.
 */
export const createMockUsers = ( count: number, teams?: number ) => {
  const mockUsers = [] as IUser[];

  for ( let i = 0; i < count; i++ ) {
    const user = createUser( teams );

    mockUsers.push( user );
  }

  return mockUsers;
};
