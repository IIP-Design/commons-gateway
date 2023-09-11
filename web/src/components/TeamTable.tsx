// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useMemo, useState } from 'react';
import type { FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import type { ColumnDef } from '@tanstack/react-table';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { buildQuery } from '../utils/api';
import { Table, defaultColumnDef } from './Table';
import { TeamModal } from './TeamModal';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import style from '../styles/table.module.scss';
import btnStyle from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const TeamTable: FC = () => {
  const [teams, setTeams] = useState<ITeam[]>( [] );

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

  const columns = useMemo<ColumnDef<ITeam>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => (
          <TeamModal
            anchor={<span style={{cursor:'pointer'}}>{info.getValue() as string}</span>}
            team={info.row.original}
            setTeams={setTeams}
          />
        ),
      },
      {
        ...defaultColumnDef( 'active' ),
        cell: info => {
          const isActive = info.getValue() as boolean;
          return (
            <span className={ style.status }>
              <span className={ isActive ? style.active : style.inactive } />
              { isActive ? 'Active' : 'Inactive' }
            </span>
          );
        },
      },
    ],
    [teams]
  );

  return (
    <div className={ style.container }>
      <TeamModal
        anchor={
          <button
            className={ `${style['add-btn']} ${btnStyle.btn}` }
            type="button"
          >
            + New Team
          </button>
        }
        setTeams={setTeams}
      />
      <Table
        {
          ...{
            data: teams,
            columns,
          }
        }
      />
    </div>
  );
}

export default TeamTable;
