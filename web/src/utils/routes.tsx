import type { TUserRole } from "../stores/current-user";

interface IRoute {
  href: string;
  name?: string;
  rolesAccessible?: TUserRole[];
}

export const routes: IRoute[] = [
  {
    href: '',
    name: 'Home',
  },
  {
    href: 'teams',
    rolesAccessible: [ 'admin', 'super admin' ],
  },
  {
    href: 'upload',
  },
]


export const filterRoutes = ( userRole: TUserRole ) => {
  return routes.filter( r => !r.rolesAccessible || r.rolesAccessible.includes( userRole ) );
}
