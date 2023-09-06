// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useEffect, useMemo, useState } from 'react';
import type { FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import {
  Column,
  Table as ReactTable,
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  ColumnDef,
  flexRender,
  getSortedRowModel,
  SortingState,
} from '@tanstack/react-table'

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { buildQuery } from '../utils/api';
import currentUser from '../stores/current-user';
import type { TUserRole } from '../stores/current-user';
import { isGuestActive } from '../utils/guest';
import { getTeamName } from '../utils/team';
import { selectSlice } from '../utils/arrays';
import { userIsSuperAdmin } from '../utils/auth';
import { renderCountWidget, setIntermediatePagination } from '../utils/pagination';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import style from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
interface IUserEntry {
  email: string;
  expiration: string;
  familyName: string;
  givenName: string;
  role: TUserRole;
  team: string;
}

interface ITeam {
  active: boolean;
  id: string;
  name: string;
}

interface ITableProps {
  data: IUserEntry[];
  columns: ColumnDef<IUserEntry>[];
}

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////

const UserTable: FC = () => {
  // Set the high and low ends of the view toggle.
  const LOW_VIEW = 30;
  const HIGH_VIEW = 90;

  const [viewCount, setViewCount] = useState( LOW_VIEW );
  const [viewOffset, setViewOffset] = useState( 0 );
  const [userList, setUserList] = useState<IUserEntry[]>( selectSlice( [], viewCount, viewOffset ) );
  const [userCount, setUserCount] = useState( userList.length ); // eslint-disable-line no-unused-vars, @typescript-eslint/no-unused-vars
  const [teams, setTeams] = useState<ITeam[]>( [] );

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

const PaginationFooter = <T, >( { table }: { table: ReactTable<T> } ) => {
  return (
    <div className="flex items-center gap-2">
        <button
          className="border rounded p-1"
          onClick={() => table.setPageIndex(0)}
          disabled={!table.getCanPreviousPage()}
        >
          {'<<'}
        </button>
        <button
          className="border rounded p-1"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          {'<'}
        </button>
        <button
          className="border rounded p-1"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          {'>'}
        </button>
        <button
          className="border rounded p-1"
          onClick={() => table.setPageIndex(table.getPageCount() - 1)}
          disabled={!table.getCanNextPage()}
        >
          {'>>'}
        </button>
        <span className="flex items-center gap-1">
          <div>Page</div>
          <strong>
            {table.getState().pagination.pageIndex + 1} of{' '}
            {table.getPageCount()}
          </strong>
        </span>
        <span className="flex items-center gap-1">
          | Go to page:
          <input
            type="number"
            defaultValue={table.getState().pagination.pageIndex + 1}
            onChange={e => {
              const page = e.target.value ? Number(e.target.value) - 1 : 0
              table.setPageIndex(page)
            }}
            className="border p-1 rounded w-16"
          />
        </span>
        <select
          value={table.getState().pagination.pageSize}
          onChange={e => {
            table.setPageSize(Number(e.target.value))
          }}
        >
          {[10, 20, 30, 40, 50].map(pageSize => (
            <option key={pageSize} value={pageSize}>
              Show {pageSize}
            </option>
          ))}
        </select>
      </div>
  );
}

const Filter = ({
  column,
  table,
}: {
  column: Column<any, any>
  table: ReactTable<any>
} ) => {
  const firstValue = table
    .getPreFilteredRowModel()
    .flatRows[0]?.getValue(column.id)

  const columnFilterValue = column.getFilterValue()

  return typeof firstValue === 'number' ? (
    <div className="flex space-x-2">
      <input
        type="number"
        value={(columnFilterValue as [number, number])?.[0] ?? ''}
        onChange={e =>
          column.setFilterValue((old: [number, number]) => [
            e.target.value,
            old?.[1],
          ])
        }
        placeholder={`Min`}
        className="w-24 border shadow rounded"
      />
      <input
        type="number"
        value={(columnFilterValue as [number, number])?.[1] ?? ''}
        onChange={e =>
          column.setFilterValue((old: [number, number]) => [
            old?.[0],
            e.target.value,
          ])
        }
        placeholder={`Max`}
        className="w-24 border shadow rounded"
      />
    </div>
  ) : (
    <input
      type="text"
      value={(columnFilterValue ?? '') as string}
      onChange={e => column.setFilterValue(e.target.value)}
      placeholder={`Search...`}
      className="w-36 border shadow rounded"
    />
  )
}

const Table: FC<ITableProps> = ( { data, columns }: ITableProps ) => {
  const [sorting, setSorting] = useState<SortingState>([]);

  const table = useReactTable<IUserEntry>({
    data,
    columns,
    state: {
      sorting,
    },

    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
    <div className={ style.container }>
      <table className={ `${style.table} ${style['user-table']}` }>
        <thead>
          {table.getHeaderGroups().map(headerGroup => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map(header => {
                return (
                  <th key={header.id} colSpan={header.colSpan}>
                    {header.isPlaceholder ? null : (
                      <div
                      {...{
                        className: header.column.getCanSort()
                          ? 'cursor-pointer select-none'
                          : '',
                        onClick: header.column.getToggleSortingHandler(),
                      }}
                      >
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                        {{
                          asc: ' 🔺',
                          desc: ' 🔻',
                        }[header.column.getIsSorted() as string] ?? null}
                        {header.column.getCanFilter() ? (
                          <div>
                            <Filter column={header.column} table={table} />
                          </div>
                        ) : null}
                      </div>
                    )}
                  </th>
                )
              })}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map(row => {
            return (
              <tr key={row.id}>
                {row.getVisibleCells().map(cell => {
                  return (
                    <td key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </td>
                  )
                })}
              </tr>
            )
          })}
        </tbody>
      </table>
      { table.getPageCount() > 1 && <PaginationFooter table={table} /> }
    </div>
  );
}

function capitalize(str: string): string {
  if( !str || !str.length ) {
      return "";
  }

  const lower = str.toLowerCase();
  return lower.substring( 0, 1 ).toUpperCase() + lower.substring( 1, lower.length );
}

const titleCase = ( str: string ) => {
  const parts =
        str
            ?.replace(/([A-Z])+/g, capitalize)
            // eslint-disable-next-line no-useless-escape
            ?.split(/(?=[A-Z])|[\.\-\s_]/)
            .map(x => x.toLowerCase()) ?? [];

    if( parts.length === 0 ) {
        return "";
    }

    parts[0] = capitalize(parts[0]);

    return parts.reduce( ( acc, part ) => {
        return `${acc} ${part.charAt(0).toUpperCase()}${part.slice(1)}`;
    } );
}

const defaultColumnDef = <T, >( key: keyof T ): ColumnDef<T> => {
  return {
    accessorFn: row => row[key],
    id: key as string,
    cell: info => info.getValue(),
    header: () => <span>{titleCase(key as string)}</span>,
    footer: props => props.column.id,
    enableSorting: true,
  }
};

const UserTable2: FC = () => {
  const [users, setUsers] = useState<IUserEntry[]>( [] );
  const [teams, setTeams] = useState<ITeam[]>( [] );

  useEffect( () => {
    const body = userIsSuperAdmin() ? {} : { team: currentUser.get().team };

    const getTeams = async () => {
      const response = await buildQuery( 'teams', null, 'GET' );
      const { data } = await response.json();

      if ( data ) {
        setTeams( data );
      }
    };

    const getUsers = async () => {
      const response = await buildQuery( 'guests', body );
      const { data } = await response.json();

      if ( data ) {
        setUsers( data );
      }
    };

    getTeams().then( () => getUsers() );
  }, [] );

  const columns = useMemo<ColumnDef<IUserEntry>[]>(
    () => [
      defaultColumnDef( 'givenName' ),
      defaultColumnDef( 'familyName' ),
      defaultColumnDef( 'email' ),
      {
        ...defaultColumnDef( 'team' ),
        cell: info => getTeamName( info.getValue() as string, teams ),
      },
      {
        ...defaultColumnDef( 'expiration' ),
        cell: info => {
          const exp = info.getValue() as string;
          const isActive = isGuestActive( exp );
          return (
            <span className={ style.status }>
              <span className={ isActive ? style.active : style.inactive } />
              { isActive ? 'Active' : 'Inactive' }
            </span>
          );
        },
      },
    ],
    []
  );

  return (
    <div className={ style.container }>
      <Table
        {
          ...{ data: users, columns }
        }
      />
    </div>
  );
}

export default UserTable2;
