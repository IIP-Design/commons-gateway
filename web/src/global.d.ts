declare module '*.scss';

type Nullable<T> = T | null;

interface ITeam {
  id: string
  name: string
  aprimoName: string
  active: boolean
}

interface IUser {
  firstName: string
  lastName: string
  email: string
  team: string
  active: boolean
}
