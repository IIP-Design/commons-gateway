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
    rolesAccessible: [ 'admin', 'super admin' ]
  },
  {
    href: 'teams',
    rolesAccessible: [ 'super admin' ],
  },
  {
    href: 'upload',
  },
]


export const filterRoutes = ( userRole: TUserRole ) => {
  return routes.filter( r => !r.rolesAccessible || r.rolesAccessible.includes( userRole ) );
}
