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
    href: 'pendingInvites',
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
    href: 'uploaderUsers',
    name: 'Uploaders',
    rolesAccessible: ['guest admin'],
  },
  {
    href: 'upload',
  },
];

export const filterRoutes = ( userRole: TUserRole ) => routes
  .filter( r => !r.rolesAccessible || r.rolesAccessible.includes( userRole ) );
