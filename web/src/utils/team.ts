/**
 * Retrieves the name of a team based on it's id.
 * @param id The team id in questions.
 * @param teams A list of teams.
 * @returns The given team name.
 */
export const getTeamName = (id: string, teams: ITeam[]) => {
  const found = teams.filter(team => team.id === id )?.[0] || {};

  return found.teamName || '';
}