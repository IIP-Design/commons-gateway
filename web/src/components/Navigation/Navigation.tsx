import currentUser from '../../stores/current-user';
import { logout } from '../../utils/login';
import { filterRoutes } from '../../utils/routes';

import './Navigation.scss';

const capitalized = ( word: string ) => word.charAt(0).toUpperCase() + word.slice(1);

export const NavBar = () => {
  const userRole = currentUser.get().role;
  const { pathname } = new URL(window.location.href);
  const currentPath = pathname.replaceAll('/', '');

  const filteredRoutes = filterRoutes( userRole || 'guest' );
  
  return (
    <nav className="nav-links">
      <ul>
        {
          filteredRoutes.map( ( { href, name } ) => (
            <li key={href}>
              <a className={currentPath === href ? 'active' : ''} href={`/${href}`}>{name || capitalized(href)}</a>
            </li>
          ) )
        }
        <li className="sign-out">
          <button id="sign-out-hamburger" type="button" onClick={logout}>Sign Out</button>
        </li>
      </ul>
    </nav>
  );
}

export default NavBar;