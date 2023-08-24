declare module '*.scss';

interface ITeam {
  id: string
  name: string
  active: boolean
}

interface IUser {
  firstName: string
  lastName: string
  email: string
  team: string
  active: boolean
}
