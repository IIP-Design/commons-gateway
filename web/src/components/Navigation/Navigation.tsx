import currentUser from '../../stores/current-user';
import { logout } from '../../utils/login';
import { filterRoutes } from '../../utils/routes';

import './Navigation.scss';

const capitalized = ( word: string ) => word.charAt( 0 ).toUpperCase() + word.slice( 1 );

export const Navigation = () => {
  const userRole = currentUser.get().role;
  const { pathname } = new URL( window.location.href );
  const currentPath = pathname.replaceAll( '/', '' );

  const filteredRoutes = filterRoutes( userRole || 'guest' );

  return (
    <nav className="nav-links">
      <ul id="nav-list">
        {
          filteredRoutes.map( ( { href, name, rolesAccessible } ) => (
            <li key={ href } className="nav-li" data-roles={ rolesAccessible?.join( '-' ) || '' }>
              <a className={ `${currentPath === href ? 'active' : ''} filterable-link` } href={ `/${href}` }>
                { name || capitalized( href ) }
              </a>
            </li>
          ) )
        }
        <li className="sign-out">
          <button id="sign-out-hamburger" type="button" onClick={ logout }>Sign Out</button>
        </li>
      </ul>
    </nav>
  );
};

export default Navigation;
