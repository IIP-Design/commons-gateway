import type { TUserRole } from './types';

interface IRoute {
  href: string;
  name?: string;
  rolesAccessible?: TUserRole[];
}

export const routes: IRoute[] = [
  {
    href: '',
    name: 'Partners',
    rolesAccessible: ['admin', 'super admin'],
  },
  {
    href: 'pending-invites',
    name: 'Invites',
    rolesAccessible: ['admin', 'super admin'],
  },
  {
    href: 'admins',
    rolesAccessible: ['super admin'],
  },
  {
    href: 'teams',
    rolesAccessible: ['super admin'],
  },
  {
    href: 'uploader-users',
    name: 'My Team',
    rolesAccessible: ['guest admin'],
  },
  {
    href: 'upload',
  },
  {
    href: 'profile',
    name: 'Profile',
    rolesAccessible: ['guest admin', 'guest'],
  },
];

export const filterRoutes = ( userRole: TUserRole ) => routes
  .filter( r => !r.rolesAccessible || r.rolesAccessible.includes( userRole ) );
