declare module '*.scss';

interface ITeam {
  id: string
  teamName: string
}

interface IUser {
  firstName: string
  lastName: string
  email: string
  team: string
  active: boolean
}
