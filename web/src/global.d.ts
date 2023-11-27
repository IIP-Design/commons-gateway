declare module '*.scss';

type Nullable<T> = T | null;
type ValueOf<T> = T[keyof T];

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

interface IMfaRequest {
  id: string
  code: string
}
