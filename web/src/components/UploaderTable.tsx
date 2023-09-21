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
import currentUser from '../stores/current-user';
import type { IUserEntry, WithUiData } from '../utils/types';
import { buildQuery } from '../utils/api';
import { Table, defaultColumnDef } from './Table';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import style from '../styles/table.module.scss';
import { daysUntil } from '../utils/dates';
import { isGuestActive } from '../utils/guest';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IUploader extends IUserEntry {
  dateInvited: string;
  proposer: string;
  inviter: string;
  pending: boolean;
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const UploaderTable: FC = () => {
  const [users, setUsers] = useState<WithUiData<IUploader>[]>( [] );

  useEffect( () => {
    const body = { team: currentUser.get().team };

    const getUsers = async () => {
      const response = await buildQuery( 'guests/uploaders', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers(
          data.map( ( user: IUploader ) => {
            return {
              ...user,
              name: `${user.givenName} ${user.familyName}`,
              active: isGuestActive( user.expiration ),
            };
          } )
        );
      }
    };

    getUsers();
  }, [] );

  const columns = useMemo<ColumnDef<WithUiData<IUploader>>[]>(
    () => [
      {
        ...defaultColumnDef( 'name' ),
        cell: info => <a href={`/editUser?id=${info.row.getValue('email')}`}>{info.getValue() as string}</a>,
      },
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'proposer' ),
        cell: info => {
          const { String, Valid } = info.getValue() as any;
          return Valid ? String : null;
        }
      },
      {
        ...defaultColumnDef( 'inviter' ),
        cell: info => {
          const { String, Valid } = info.getValue() as any;
          return Valid ? String : null;
        }
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
      {
        ...defaultColumnDef( 'pending' ),
        header: 'Status',
        cell: info => {
          const isPending = info.getValue() as boolean;
          return (
            <span className={ style.status }>
              <span className={ isPending ? style.inactive : style.active } />
              { isPending ? 'Pending' : 'Approved' }
            </span>
          );
        },
      },
    ],
    []
  );

  return (
    <div className={ style.container }>
      { users.length ?
        <Table
          {
            ...{
              data: users,
              columns,
              additionalTableClasses: [ 'user-table' ],
            }
          }
        />
        : <p>No data to show</p>
      }
    </div>
  );
}

export default UploaderTable;
