// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useState } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import {
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  flexRender,
  getSortedRowModel,
} from '@tanstack/react-table';
import type { Column, Table as ReactTable, ColumnDef, SortingState, Header } from '@tanstack/react-table';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { titleCase } from '../utils/string';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import tableStyles from '../styles/table.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Types and Interfaces
// ////////////////////////////////////////////////////////////////////////////
export type ITableProps<DataType> = {
  readonly data: DataType[];
  readonly columns: ColumnDef<DataType>[];
  readonly additionalTableClasses?: string[];
}

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
export const defaultColumnDef = <T, >( key: keyof T, header?: string ): ColumnDef<T> => ( {
  accessorFn: row => row[key],
  id: key as string,
  cell: info => info.getValue(),
  header: () => <span>{ header || titleCase( key as string ) }</span>,
  footer: props => props.column.id,
  enableSorting: true,
} );

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const Filter = ( {
  column,
  table,
}: {
  column: Column<any, any>  // eslint-disable-line @typescript-eslint/no-explicit-any
  table: ReactTable<any>    // eslint-disable-line @typescript-eslint/no-explicit-any
} ) => {
  const skipFiltering = ( column.id.startsWith( '_' ) );

  const firstValue = table
    .getPreFilteredRowModel()
    .flatRows[0]?.getValue( column.id );
  const colType = ( skipFiltering ? null : typeof firstValue );
  const columnFilterValue = column.getFilterValue();

  let ret: JSX.Element | null = null;

  switch ( colType ) {
    case 'number': {
      ret = (
        <div>
          <input
            type="number"
            value={ ( columnFilterValue as [number, number] )?.[0] ?? '' }
            onChange={ e => column.setFilterValue( ( old: [number, number] ) => [
              e.target.value,
              old?.[1],
            ] ) }
            placeholder="Min"
            aria-label={ `${titleCase( column.id )} minimum value` }
          />
          <input
            type="number"
            value={ ( columnFilterValue as [number, number] )?.[1] ?? '' }
            onChange={ e => column.setFilterValue( ( old: [number, number] ) => [
              old?.[0],
              e.target.value,
            ] ) }
            placeholder="Max"
            aria-label={ `${titleCase( column.id )} maximum value` }
          />
        </div>
      );
      break;
    }
    case 'string': {
      ret = (
        <input
          aria-label={ `Search ${titleCase( column.id )}` }
          id={ `search-${column.id}` }
          placeholder="Search..."
          style={ { width: '100%' } }
          type="text"
          value={ ( columnFilterValue ?? '' ) as string }
          onChange={ e => column.setFilterValue( e.target.value ) }
        />
      );
      break;
    }
    default:
      break;
  }

  return ret;
};

export const PaginationFooter = <T, >( { table }: { readonly table: ReactTable<T> } ) => {
  const revNavDisabled = !table.getCanPreviousPage();
  const revNavStyle = revNavDisabled ? { cursor: 'default' } : {};

  const fwdNavDisabled = !table.getCanNextPage();
  const fwdNavStyle = fwdNavDisabled ? { cursor: 'default' } : {};

  return (
    <div className={ tableStyles['pagination-footer'] }>
      <button
        className={ tableStyles['pagination-button'] }
        style={ revNavStyle }
        onClick={ () => table.setPageIndex( 0 ) }
        disabled={ revNavDisabled }
        type="button"
      >
        { '<<' }
      </button>
      <button
        className={ tableStyles['pagination-button'] }
        style={ revNavStyle }
        onClick={ () => table.previousPage() }
        disabled={ revNavDisabled }
        type="button"
      >
        { '<' }
      </button>
      <span className={ tableStyles['go-to-page-container'] }>
        <div style={ { display: 'inline' } }>Page</div>
        <strong>
          { `${table.getState().pagination.pageIndex + 1} of ${table.getPageCount()}` }
        </strong>
      </span>
      <button
        className={ tableStyles['pagination-button'] }
        style={ fwdNavStyle }
        onClick={ () => table.nextPage() }
        disabled={ fwdNavDisabled }
        type="button"
      >
        { '>' }
      </button>
      <button
        className={ tableStyles['pagination-button'] }
        style={ fwdNavStyle }
        onClick={ () => table.setPageIndex( table.getPageCount() - 1 ) }
        disabled={ fwdNavDisabled }
        type="button"
      >
        { '>>' }
      </button>
      <span className={ tableStyles['go-to-page-container'] }>
        { 'Go to page: ' }
        <input
          aria-label="Go to page"
          className={ tableStyles['page-num-select'] }
          id="goto-page-number-input"
          defaultValue={ table.getState().pagination.pageIndex + 1 }
          max={ table.getPageCount() }
          min={ 1 }
          type="number"
          onChange={ e => {
            const page = e.target.value ? Number( e.target.value ) - 1 : 0;

            table.setPageIndex( page );
          } }
        />
      </span>
      <span className={ tableStyles['go-to-page-container'] }>
        { 'Show ' }
        <select
          aria-label="Select number of results per page"
          className={ tableStyles['page-show-select'] }
          id="show-page-number-select"
          style={ { width: 'auto', marginLeft: '0.375em', marginRight: '0.375em' } }
          value={ table.getState().pagination.pageSize }
          onChange={ e => {
            table.setPageSize( Number( e.target.value ) );
          } }
        >
          {
            [
              10, 20, 30, 40, 50,
            ].filter( ( _, idx, arr ) => idx === 0 || arr[idx - 1] < table.getFilteredRowModel().rows.length )
              .map( pageSize => (
                <option key={ pageSize } value={ pageSize }>
                  { pageSize }
                </option>
              ) )
          }
        </select>
        { ' per page' }
      </span>
    </div>
  );
};

export const Table = <DataType, >( { data, columns, additionalTableClasses }: ITableProps<DataType> ) => {
  const [sorting, setSorting] = useState<SortingState>( [] );

  const makeHeaderKeyDownHandler = ( header: Header<DataType, unknown> ) => {
    const handler = header.column.getToggleSortingHandler();

    return ( e: React.KeyboardEvent ) => {
      if ( e.key === 'Enter' && handler ) {
        handler( e );
      }
    };
  };

  const table = useReactTable<DataType>( {
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
  } );

  const tableClasses = ( additionalTableClasses || [] ).map( c => tableStyles[c] ).filter( Boolean )
    .join( ' ' );

  return (
    <div className={ tableStyles.container }>
      <table className={ `${tableStyles.table} ${tableClasses}` }>
        <thead>
          { table.getHeaderGroups().map( headerGroup => (
            <tr key={ headerGroup.id }>
              { headerGroup.headers.map( header => (
                <th key={ header.id } colSpan={ header.colSpan } style={ { verticalAlign: 'top' } }>
                  { header.isPlaceholder
                    ? null
                    : (
                      <div
                        className={ header.column.getCanSort() ? tableStyles['sortable-header'] : '' }
                      >
                        <span
                          onClick={ header.column.getToggleSortingHandler() }
                          onKeyDown={ makeHeaderKeyDownHandler( header ) }
                          role="button"
                          tabIndex={ 0 }
                        >
                          { flexRender(
                            header.column.columnDef.header,
                            header.getContext(),
                          ) }
                        </span>
                        { {
                          asc: ' 🔺',
                          desc: ' 🔻',
                        }[header.column.getIsSorted() as string] ?? null }
                        { header.column.getCanFilter()
                          ? (
                            <div>
                              <Filter column={ header.column } table={ table } />
                            </div>
                          )
                          : null }
                      </div>
                    ) }
                </th>
              ) ) }
            </tr>
          ) ) }
        </thead>
        <tbody>
          { table.getRowModel().rows.map( row => (
            <tr key={ row.id }>
              { row.getVisibleCells().map( cell => (
                <td key={ cell.id }>
                  { flexRender(
                    cell.column.columnDef.cell,
                    cell.getContext(),
                  ) }
                </td>
              ) ) }
            </tr>
          ) ) }
        </tbody>
      </table>
      <hr />
      <PaginationFooter table={ table } />
    </div>
  );
};
