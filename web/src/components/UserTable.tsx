import { useEffect, useState } from 'react';
import type { FC } from 'react';

import { buildQuery } from '../utils/api';
import currentUser from '../stores/current-user';
import { isGuestActive } from '../utils/guest';
import { getTeamName } from '../utils/team';
import { selectSlice } from '../utils/arrays';
import { userIsSuperAdmin } from '../utils/auth';
import { renderCountWidget, setIntermediatePagination } from '../utils/pagination';

import style from '../styles/table.module.scss';

const UserTable: FC = () => {
  // Set the high and low ends of the view toggle.
  const LOW_VIEW = 30;
  const HIGH_VIEW = 90;

  const [viewCount, setViewCount] = useState( LOW_VIEW );
  const [viewOffset, setViewOffset] = useState( 0 );
  const [userList, setUserList] = useState( selectSlice( [], viewCount, viewOffset ) );
  const [userCount, setUserCount] = useState( userList.length ); // eslint-disable-line no-unused-vars, @typescript-eslint/no-unused-vars
  const [teams, setTeams] = useState( [] );

  useEffect( () => {
    const body = userIsSuperAdmin() ? {} : { team: currentUser.get().team };

    const getUsers = async () => {
      const response = await buildQuery( 'guests', body );
      const { data } = await response.json();

      if ( data ) {
        setUserList( selectSlice( data, viewCount, viewOffset ) );
        setUserCount( data.length );
      }
    };

    getUsers();
  }, [viewCount, viewOffset] );

  useEffect( () => {
    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setTeams( data );
      }
    };

    getTeams();
  }, [] );

  // How many more users are left to the end of the list.
  const remainingScroll = userCount - ( viewCount * viewOffset );

  /**
   * Advance the table scroll forward or backwards.
   * @param dir The direction of scroll, positive for forward, negative for back.
   */
  const turnPage = ( dir: 1 | -1 ) => {
    setViewOffset( viewOffset + dir );
  };

  /**
   * Advance the table scroll to a give page of results.
   * @param page The page to navigate to.
   */
  const goToPage = ( page: number ) => {
    setViewOffset( page - 1 ); // Adjustment since offsets start at zero.
  };

  /**
   * Toggle the number of items displayed in the table.
   * @param count How many to show.
   */
  const changeViewCount = ( count: number ) => {
    setViewCount( count );
    // We reset the offset in case the current
    // offset * new count is more than total users
    setViewOffset( 0 );
  };

  return (
    <div className={ style.container }>
      <div>
        <div className={ style.controls }>
          <span>{ renderCountWidget( userCount, viewCount, viewOffset ) }</span>
          { userCount > LOW_VIEW && (
            <div className={ style.count }>
              <span>View:</span>
              <button
                className={ style['pagination-btn'] }
                onClick={ () => changeViewCount( LOW_VIEW ) }
                disabled={ viewCount === LOW_VIEW }
                type="button"
              >
                { LOW_VIEW }
              </button>
              <span>|</span>
              <button
                className={ style['pagination-btn'] }
                onClick={ () => changeViewCount( HIGH_VIEW ) }
                disabled={ viewCount === HIGH_VIEW }
                type="button"
              >
                { HIGH_VIEW }
              </button>
            </div>
          ) }
        </div>
        <table className={ `${style.table} ${style['user-table']}` }>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Team Name</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            { userList && ( userList.map( user => (
              <tr key={ user.email }>
                <td>
                  <a href={ `/editUser?id=${user.email}` } style={ { padding: '0.3rem' } }>
                    { `${user.givenName} ${user.familyName}` }
                  </a>
                </td>
                <td>{ user.email }</td>
                <td>{ getTeamName( user.team, teams ) }</td>
                <td className={ style.status }>
                  <span className={ isGuestActive( user.expiration ) ? style.active : style.inactive } />
                  { isGuestActive( user.expiration ) ? 'Active' : 'Inactive' }
                </td>
              </tr>
            ) ) ) }
          </tbody>
        </table>
      </div>
      { viewCount < userCount && (
        <div className={ style.pagination }>
          <button
            className={ style['pagination-btn'] }
            type="button"
            onClick={ () => turnPage( -1 ) }
            disabled={ viewOffset < 1 }
          >
            { '< Prev' }
          </button>
          { setIntermediatePagination( userCount, viewCount ).length >= 3 && (
            <span className={ style['pagination-intermediate'] }>
              { setIntermediatePagination( userCount, viewCount ).map( page => (
                <button
                  key={ page }
                  className={ style['pagination-btn'] }
                  disabled={ viewOffset + 1 === page }
                  onClick={ () => goToPage( page ) }
                  type="button"
                >
                  { page }
                </button>
              ) ) }
            </span>
          ) }
          <button
            className={ style['pagination-btn'] }
            type="button"
            onClick={ () => turnPage( 1 ) }
            disabled={ remainingScroll <= viewCount }
          >
            { 'Next >' }
          </button>
        </div>
      ) }
    </div>
  );
};

export default UserTable;
