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
import { defaultColumnDef } from './Table';
import { TeamModal } from './TeamModal/TeamModal';
import TableWrapper from './TableWrapper';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import tableStyles from '../styles/table.module.scss';
import btnStyle from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const TeamTable: FC = () => {
  const [teams, setTeams] = useState<ITeam[]>( [] );
  const [loading, setLoading] = useState<boolean>( true );

  useEffect( () => {
    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setTeams( data );
      }
    };

    getTeams().finally( () => setLoading( false ) );
  }, [] );

  const columns = useMemo<ColumnDef<ITeam>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => (
          <TeamModal
            anchor={ <span>{ info.getValue() as string }</span> }
            team={ info.row.original }
            setTeams={ setTeams }
          />
        ),
      },
      {
        ...defaultColumnDef( 'active' ),
        cell: info => {
          const isActive = info.getValue() as boolean;

          return (
            <span className={ tableStyles.status }>
              <span className={ isActive ? tableStyles.active : tableStyles.inactive } />
              { isActive ? 'Active' : 'Inactive' }
            </span>
          );
        },
      },
    ],
    [],
  );

  return (
    <div style={ { display: 'flex', flexDirection: 'column', marginBottom: '0.75em' } }>
      <TeamModal
        anchor={ (
          <span
            className={ `${tableStyles['add-btn']} ${btnStyle.btn}` }
            style={ { fontSize: 'var(--fontSizeSmall)' } }
          >
            + New Team
          </span>
        ) }
        setTeams={ setTeams }
      />
      <TableWrapper
        loading={ loading }
        table={ { data: teams, columns } }
      />
    </div>
  );
};

export default TeamTable;
