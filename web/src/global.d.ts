declare module '*.scss';

interface ITeam {
  id: string
  teamName: string
  active: boolean
}

interface IUser {
  firstName: string
  lastName: string
  email: string
  team: string
  active: boolean
}
