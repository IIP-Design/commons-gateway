import { useState } from 'react';
import type { FC } from 'react';

import UserTable from '../UserTable';

import style from './ExternalPartnerTables.module.scss';

const ExternalPartnerTables: FC = () => {
  const [table, setTable] = useState( 'partners' );

  return (
    <div>
      <select aria-label="User Type" className={ style.select } id="user-type-select" value={ table } onChange={ e => setTable( e.target.value ) }>
        <option value="partners">External Partners</option>
        <option value="leads">External Team Leads</option>
      </select>
      { table === 'partners' && ( <UserTable role="guest" /> ) }
      { table === 'leads' && ( <UserTable role="guest admin" /> ) }
    </div>
  );
};

export default ExternalPartnerTables;
