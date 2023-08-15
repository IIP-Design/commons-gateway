import { generateName } from "./utils";

/**
 * Creates a mock team by creating a random name and
 * assigning it the provided id value.
 * @param id 
 */
const createTeam = (id: number): ITeam => ({
  id: id.toString(),
  teamName: `Team ${generateName(5)}`,
});

/**
 * Generates a list of teams of a specified length.
 * @param count The number of teams to create.
 */
export const createMockTeams = (count: number) => {
  const mockTeams = [] as ITeam[];

  for (let i = 0; i < count; i++) {
    const team = createTeam(i + 1);
    mockTeams.push(team)
  }

  return mockTeams;
}