import type { UserRole } from "../stores/current-user";

interface IRoute {
  href: string;
  name?: string;
  rolesAccessible?: UserRole[];
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


export const filterRoutes = ( userRole: UserRole ) => {
  return routes.filter( r => !r.rolesAccessible || r.rolesAccessible.includes( userRole ) );
}
