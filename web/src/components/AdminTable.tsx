import { useEffect, useState } from 'react';
import type { FC } from 'react';

import { buildQuery } from '../utils/api';
import { getTeamName } from '../utils/team';
import { selectSlice } from '../utils/arrays';
import { renderCountWidget, setIntermediatePagination } from '../utils/pagination';

import style from '../styles/table.module.scss';

const AdminTable: FC = () => {
  // Set the high and low ends of the view toggle.
  const LOW_VIEW = 30;
  const HIGH_VIEW = 90;

  const [viewCount, setViewCount] = useState( LOW_VIEW );
  const [viewOffset, setViewOffset] = useState( 0 );
  const [adminList, setAdminList] = useState( selectSlice( [], viewCount, viewOffset ) );
  const [adminCount, setAdminCount] = useState( adminList.length ); // eslint-disable-line no-unused-vars, @typescript-eslint/no-unused-vars
  const [teams, setTeams] = useState( [] );

  useEffect( () => {
    const getAdmins = async () => {
      const response = await buildQuery( 'admins', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setAdminList( selectSlice( data, viewCount, viewOffset ) );
        setAdminCount( data.length );
      }
    };

    getAdmins();
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
  const remainingScroll = adminCount - ( viewCount * viewOffset );

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
          <span>{ renderCountWidget( adminCount, viewCount, viewOffset ) }</span>
          { adminCount > LOW_VIEW && (
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
              <th>Team</th>
              <th>Role</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            { adminList && ( adminList.map( admin => (
              <tr key={ admin.email }>
                <td>
                  <a href={ `/editAdmin?id=${admin.email}` } style={ { padding: '0.3rem' } }>
                    { `${admin.givenName} ${admin.familyName}` }
                  </a>
                </td>
                <td>{ admin.email }</td>
                <td>{ getTeamName( admin.team, teams ) }</td>
                <td>{ admin.role }</td>
                <td className={ style.status }>
                  <span className={ admin.active ? style.active : style.inactive } />
                  { admin.active ? 'Active' : 'Inactive' }
                </td>
              </tr>
            ) ) ) }
          </tbody>
        </table>
      </div>
      { viewCount < adminCount && (
        <div className={ style.pagination }>
          <button
            className={ style['pagination-btn'] }
            type="button"
            onClick={ () => turnPage( -1 ) }
            disabled={ viewOffset < 1 }
          >
            { '< Prev' }
          </button>
          { setIntermediatePagination( adminCount, viewCount ).length >= 3 && (
            <span className={ style['pagination-intermediate'] }>
              { setIntermediatePagination( adminCount, viewCount ).map( page => (
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

export default AdminTable;
