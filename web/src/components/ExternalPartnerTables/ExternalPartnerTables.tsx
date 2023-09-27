import { useState } from 'react';
import type { FC } from 'react';

import UserTable from '../UserTable';

import style from './ExternalPartnerTables.module.scss';

const ExternalPartnerTables: FC = () => {
  const [table, setTable] = useState( 'partners' );

  return (
    <div>
      <select className={ style.select } value={ table } onChange={ e => setTable( e.target.value ) } aria-label="User Type">
        <option value="partners">External Partners</option>
        <option value="leads">External Team Leads</option>
      </select>
      { table === 'partners' && ( <UserTable role="guest" /> ) }
      { table === 'leads' && ( <UserTable role="guest admin" /> ) }
    </div>
  );
};

export default ExternalPartnerTables;
