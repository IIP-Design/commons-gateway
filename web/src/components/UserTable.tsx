import { useEffect, useState } from 'react';
import type { FC } from 'react';

import { getTeamName } from '../utils/team';
import { selectSlice } from '../utils/arrays';

import style from './UserTable.module.scss';

interface IUserTableProps {
  users: IUser[]
  teams: ITeam[]
}

/**
 * Determines what to display in the scroll controls range.
 * @param total Count of all users.
 * @param show Number of users to show in the table.
 * @param offset The depth of scroll into the total (by increments of the viewCount).
 * @returns The string that makes up the top of the range.
 */
const setUpperBound = (total: number, show: number, offset: number) => {
  const remaining = total - (show * offset )
  const theoreticalTop = show * (offset + 1);

  if ( remaining === 1 ) {
    return '';
  }
  
  if ( total < show || remaining < show ) {
    return `-${total}`
  }

  return `-${theoreticalTop}`;
}

/**
 * Generates the string showing the pagination status of the table.
 * @param userCount The total number of users.
 * @param viewCount The number of items shown at any one time.
 * @param viewOffset The depth of scroll into the total (by increments of the viewCount).
 * @returns The full count string.
 */
const renderCountWidget = (userCount: number, viewCount: number, viewOffset: number) => {
  const start = viewOffset * viewCount + 1;
  const end = setUpperBound(userCount, viewCount, viewOffset);

  return `${start}${end} out of ${userCount}`
}

/**
 * Calculates how many pages of results the table has.
 * @param total The total number of users.
 * @param show The number of items shown at any one time.
 */
const setIntermediatePagination = (total: number, show: number) => {
  const divisions = Math.floor(total/show);

  const pages = [];

  for (let i = 0; i < divisions; i++) {
    pages.push(i + 1);
  }

  return pages;
}

const UserTable: FC<IUserTableProps> = ({users, teams}) => {
  // Set the high and low ends of the view toggle.
  const LOW_VIEW = 30;
  const HIGH_VIEW = 90;

  const [userCount] = useState(users.length);
  const [viewCount, setViewCount] = useState(LOW_VIEW);
  const [viewOffset, setViewOffset] = useState(0);
  const [userList, setUserList] = useState(selectSlice(users, viewCount, viewOffset));

  useEffect(() => {
    setUserList(selectSlice(users, viewCount, viewOffset))
  }, [users, viewCount, viewOffset])

  // How many more users are left to the end of the list.
  const remainingScroll = userCount - (viewCount * viewOffset);

  /**
   * Advance the table scroll forward or backwards.
   * @param dir The direction of scroll, positive for forward, negative for back.
   */
  const turnPage = (dir: 1 | -1) => {
    setViewOffset(viewOffset + dir);
  }

  /**
   * Advance the table scroll to a give page of results.
   * @param page The page to navigate to.
   */
  const goToPage = (page: number) => {
    setViewOffset(page - 1) // Adjustment since offsets start at zero.
  }

  /**
   * Toggle the number of items displayed in the table.
   * @param count How many to show.
   */
  const changeViewCount = (count: number) => {
    setViewCount(count);
    // We reset the offset in case the current
    // offset * new count is more than total users
    setViewOffset(0);
  }

  return (
    <div className={style.container}>
      <div>
        <div className={style.controls}>
          <span>{renderCountWidget(userCount, viewCount, viewOffset)}</span>
          { userCount > LOW_VIEW && (
            <div className={style.count}>
              <span>View:</span>
              <button
                className={style['pagination-btn']}
                onClick={() => changeViewCount(LOW_VIEW)}
                disabled={viewCount === LOW_VIEW}
              >
                {LOW_VIEW}
              </button>
              <span>|</span>
              <button
                className={style['pagination-btn']}
                onClick={() => changeViewCount(HIGH_VIEW)}
                disabled={viewCount === HIGH_VIEW}
              >
                {HIGH_VIEW}
              </button>
            </div>
          )}
        </div>
        <table className={style.table}>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Team Name</th>
              <th>Account Status</th>
            </tr>
          </thead>
          <tbody>
            {userList && ( userList.map( user => (
              <tr key={user.email}>
                <td>{user.firstName} {user.lastName}</td>
                <td>{user.email}</td>
                <td>{getTeamName(user.team, teams)}</td>
                <td className={style.status}>
                  <span className={user.active ? style.active : style.inactive}/>
                  {user.active ? 'Active' : 'Inactive'}
                </td>
              </tr>
            )))}
          </tbody>
        </table>
      </div>
      {viewCount < userCount && (
        <div className={style.pagination}>
          <button
            className={style['pagination-btn']}
            type="button"
            onClick={() => turnPage(-1)}
            disabled={viewOffset < 1}
          >
            {`< Prev`}
          </button>
          { setIntermediatePagination(userCount, viewCount).length > 3 && (
            <span className={style['pagination-intermediate']}>
              {setIntermediatePagination(userCount, viewCount).map(page => (
                <button
                  key={page}
                  className={style['pagination-btn']}
                  disabled={viewOffset + 1 === page}
                  onClick={() => goToPage(page)}
                >
                  {page}
                </button>
              ))}
            </span>
          )}
          <button
            className={style['pagination-btn']}
            type="button"
            onClick={() => turnPage(1)}
            disabled={remainingScroll <= viewCount}
          >
            {`Next >`}
          </button>
        </div>
      )}
    </div>
  )
};

export default UserTable;