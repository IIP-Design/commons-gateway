import { useEffect, useState } from 'react';
import type { FC } from 'react';

import ToggleSwitch from './ToggleSwitch/ToggleSwitch';
import { buildQuery } from '../utils/api';
import { selectSlice } from '../utils/arrays';
import { renderCountWidget, setIntermediatePagination } from '../utils/pagination';

import style from '../styles/table.module.scss';
import btnStyle from '../styles/button.module.scss';

const TeamTable: FC = () => {
  // Set the high and low ends of the view toggle.
  const LOW_VIEW = 30;
  const HIGH_VIEW = 90;

  // State of the results pagination.
  const [viewCount, setViewCount] = useState( LOW_VIEW );
  const [viewOffset, setViewOffset] = useState( 0 );
  // State of the full team list.
  const [teamList, setTeamList] = useState( selectSlice( [], viewCount, viewOffset ) );
  const [teamCount, setTeamCount] = useState( teamList.length );
  // State used when add/editing a team.
  const [editing, setEditing] = useState( '' );
  const [newName, setNewName] = useState( '' );

  useEffect( () => {
    // Retrieve the full list of teams from the API.
    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setTeamList( selectSlice( data, viewCount, viewOffset ) );
        setTeamCount( data.length );
      }
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

  /**
   * Activate/deactivate a given team using the status toggle.
   * @param status The active/inactive value the team should be set to.
   * @param id The id of the team in question.
   */
  const handleStatusToggle = async ( status: boolean, id: string ) => {
    try {
      await buildQuery( 'team/update', { active: status, team: id }, 'POST' );
    } catch ( err ) {
      console.log( err );
    }
  };

  /**
   * Enable an input field to create a new team.
   */
  const addNewTeam = () => {
    // Add a placeholder team with the id `temp`
    teamList.unshift( { id: 'temp', active: true, name: '' } );
    setEditing( 'temp' );
  };

  /**
   * Enable an input field to alter an existing team name.
   * @param id The id of the team to be altered.
   * @param name The current name of the team.
   */
  const editTeam = ( id: string, name: string ) => {
    setEditing( id );
    setNewName( name );
  };

  /**
   * Abort the editing/addition of a team.
   * @param id The id of the team being altered.
   */
  const cancelEdit = ( id: string ) => {
    setEditing( '' );
    setNewName( '' );

    // If canceling a team creation, remove the placeholder team.
    if ( id === 'temp' ) {
      teamList.shift();
    }
  };

  /**
   * Sends the user inputs on the team being created/updated to the API.
   * @param team Information about the team to be created/edited.
   */
  const saveTeam = async ( team: ITeam ) => {
    let newList;

    // If the team is new, send a create request, otherwise send an update request.
    if ( team.id === 'temp' ) {
      const response = await buildQuery( 'team/create', { teamName: newName }, 'POST' );
      const { data } = await response.json();

      newList = data;
    } else {
      const response = await buildQuery( 'team/update', { active: team.active, team: team.id, teamName: newName }, 'POST' );
      const { data } = await response.json();

      newList = data;
    }

    // Update the team list with new data from the API.
    if ( newList ) {
      setTeamList( selectSlice( newList, viewCount, viewOffset ) );
      setTeamCount( newList.length );
    }

    // Reset the state used when add/editing a team.
    setEditing( '' );
    setNewName( '' );
  };

  return (
    <div className={ style.container }>
      <button
        className={ `${style['add-btn']} ${btnStyle.btn}` }
        type="button"
        onClick={ addNewTeam }
      >
        + New Team
      </button>
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
                <td>
                  { editing === team.id && (
                    <input style={ { padding: '0.3rem 0.5rem' } } type="text" value={ newName } onChange={ e => setNewName( e.target.value ) } />
                  ) }
                  { editing !== team.id && (
                    <button className={ style['pagination-btn'] } disabled={ editing !== '' } type="button" onClick={ () => editTeam( team.id, team.name ) }>
                      { team.name }
                    </button>
                  ) }
                </td>
                <td>
                  { editing === team.id && (
                    <div>
                      <button className={ btnStyle.btn } type="button" onClick={ () => saveTeam( team ) }>Save</button>
                      <button className={ btnStyle['btn-light'] } style={ { marginLeft: '1rem' } } type="button" onClick={ () => cancelEdit( team.id ) }>Cancel</button>
                    </div>
                  ) }
                  { editing !== team.id && (
                    <ToggleSwitch active={ team.active } callback={ handleStatusToggle } id={ team.id } />
                  ) }
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
