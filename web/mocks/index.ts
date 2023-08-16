import { createMockUsers } from './users';
import { createMockTeams } from './teams';

/**
 * Create a specified number of mock users and teams.
 * @param userCount The number of users to generate.
 * @param teamCount The number of teams to generate.
 */
export const mockSomeData = ( userCount: number, teamCount: number ) => {
  const teams = createMockTeams( teamCount );
  const users = createMockUsers( userCount, teamCount );

  return { teams, users };
};
