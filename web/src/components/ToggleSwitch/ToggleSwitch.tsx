import { useState } from 'react';
import type { FC } from 'react';

import style from './ToggleSwitch.module.scss';

interface IToggleSwitchProps {
  readonly id: string
  readonly active: boolean
  readonly callback: ( toggled: boolean, id: string ) => void
}

const ToggleSwitch: FC<IToggleSwitchProps> = ( { id, active, callback } ) => {
  const [toggled, setToggled] = useState( active );

  const handleToggle = () => {
    const switched = !toggled;

    setToggled( switched );
    callback( switched, id );
  };

  return (
    <label className={ style.toggle } htmlFor={ `${id}-toggle` }>
      <span style={ { marginRight: '1rem', width: '3.5rem' } }>
        { toggled ? 'Active' : 'Inactive' }
      </span>
      <input id={ `${id}-toggle` } checked={ toggled } type="checkbox" onChange={ handleToggle } />
      <span className={ style['toggle-slider'] } />
    </label>

  );
};

export default ToggleSwitch;
