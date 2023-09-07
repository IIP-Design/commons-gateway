import { useState } from 'react';
import type{ FC } from 'react';

import Modal from 'react-modal';
import { buildQuery } from '../utils/api';
import { showError } from '../utils/alert';
import ToggleSwitch from './ToggleSwitch/ToggleSwitch';

import style from '../styles/table.module.scss';
import btnStyle from '../styles/button.module.scss';

interface ITeamModalProps {
  team?: ITeam;
  setTeams: React.Dispatch<React.SetStateAction<ITeam[]>>;
  anchor: string | JSX.Element;
}

Modal.setAppElement( document.getElementById('root') as HTMLElement);

export const TeamModal: FC<ITeamModalProps> = ( { team, setTeams, anchor }: ITeamModalProps ) => {
  // Team Setup
  const [ localTeam, setTeam ] = useState<Partial<ITeam>>( team || { active: true } );
  
  // Modal Setup
  const [ modalIsOpen, setIsOpen ] = useState( false );

  // Modal Controls
  const openModal = () => setIsOpen( true );
  const closeModal = () => setIsOpen( false );

  // Update Controls
  const handleUpdate = ( key: string, value: any ) => {
    setTeam( { ...localTeam, [key]: value } );
  };

  const handleSubmit = async () => {
    const { name, id, active } = localTeam;
    if( !name ) {
      showError( 'A team must have a name' );
      return;
    }

    let newList;
    let errMessage;

    // If the team is new, send a create request, otherwise send an update request.
    if ( !id ) {

      const response = await buildQuery( 'team/create', { teamName: name }, 'POST' );
      const { data, message } = await response.json();

      newList = data;
      errMessage = message;
    } else {
      const response = await buildQuery( 'team/update', { active, team: id, teamName: name }, 'POST' );
      const { data, message } = await response.json();

      newList = data;
      errMessage = message;
    }

    // Update the team list with new data from the API.
    if ( newList ) {
      setTeams( newList );
    }

    if ( errMessage ) {
      showError( `Unable to complete your request. Reason: ${errMessage}` );
    } else {
      closeModal();
    }
  }

  return (
    <>
      <a onClick={openModal}>{anchor}</a>
      <Modal
        isOpen={modalIsOpen}
        onRequestClose={closeModal}
        contentLabel="Example Modal"
      >
        <h1>{ localTeam.id ? `Update Team ${localTeam.name}` : "Add a new Team" }</h1>
        <input
          style={ { maxWidth: '100%', padding: '0.3rem 0.5rem' } }
          type="text"
          value={ localTeam.name || '' }
          onChange={ e => handleUpdate( 'name', e.target.value ) }
          aria-label="Team Name"
        />
        { localTeam.id && <ToggleSwitch
          active={ localTeam.active ?? false }
          callback={ e => handleUpdate( 'active', e ) }
          id={ localTeam.id }
        /> }
        <button
          className={ `${style['add-btn']} ${btnStyle.btn}` }
          onClick={closeModal}
        >
          Back
        </button>
      </Modal>
    </>
  );
}