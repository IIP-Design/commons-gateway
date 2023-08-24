import { useEffect, useState } from 'react';
import type { FC } from 'react';

import ToggleSwitch from './ToggleSwitch/ToggleSwitch';
import { buildQuery } from '../utils/api';
import { selectSlice } from '../utils/arrays';
import { renderCountWidget, setIntermediatePagination } from '../utils/pagination';

import style from '../styles/table.module.scss';

const TeamTable: FC = () => {
  // Set the high and low ends of the view toggle.
  const LOW_VIEW = 30;
  const HIGH_VIEW = 90;

  const [viewCount, setViewCount] = useState( LOW_VIEW );
  const [viewOffset, setViewOffset] = useState( 0 );
  const [teamList, setTeamList] = useState( selectSlice( [], viewCount, viewOffset ) );
  const [teamCount, setTeamCount] = useState( teamList.length ); // eslint-disable-line no-unused-vars, @typescript-eslint/no-unused-vars

  useEffect( () => {
    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );

      const { data } = await response.json();

      setTeamList( selectSlice( data, viewCount, viewOffset ) );
      setTeamCount( data.length );
    };

    getTeams();
  }, [viewCount, viewOffset] );

  // How many more teams are left to the end of the list.
  const remainingScroll = teamCount - ( viewCount * viewOffset );

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
    // offset * new count is more than total teams
    setViewOffset( 0 );
  };

  const handleStatusToggle = ( status: boolean, id: string ) => {
    console.log( status, id );
  };

  return (
    <div className={ style.container }>
      <div>
        <div className={ style.controls }>
          <span>{ renderCountWidget( teamCount, viewCount, viewOffset ) }</span>
          { teamCount > LOW_VIEW && (
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
        <table className={ style.table }>
          <thead>
            <tr>
              <th>Team Name</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            { teamList && ( teamList.map( team => (
              <tr key={ team.id }>
                <td>{ team.name }</td>
                <td>
                  <ToggleSwitch active={ team.active } callback={ handleStatusToggle } id={ team.id } />
                </td>
              </tr>
            ) ) ) }
          </tbody>
        </table>
      </div>
      { viewCount < teamCount && (
        <div className={ style.pagination }>
          <button
            className={ style['pagination-btn'] }
            type="button"
            onClick={ () => turnPage( -1 ) }
            disabled={ viewOffset < 1 }
          >
            { '< Prev' }
          </button>
          { setIntermediatePagination( teamCount, viewCount ).length >= 3 && (
            <span className={ style['pagination-intermediate'] }>
              { setIntermediatePagination( teamCount, viewCount ).map( page => (
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

export default TeamTable;
