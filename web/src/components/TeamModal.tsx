// ////////////////////////////////////////////////////////////////////////////
// React Imports
// ////////////////////////////////////////////////////////////////////////////
import { useState } from 'react';
import type{ FC } from 'react';

// ////////////////////////////////////////////////////////////////////////////
// 3PP Imports
// ////////////////////////////////////////////////////////////////////////////
import Modal from 'react-modal';

// ////////////////////////////////////////////////////////////////////////////
// Local Imports
// ////////////////////////////////////////////////////////////////////////////
import { buildQuery } from '../utils/api';
import { showError } from '../utils/alert';
import ToggleSwitch from './ToggleSwitch/ToggleSwitch';

// ////////////////////////////////////////////////////////////////////////////
// Styles and CSS
// ////////////////////////////////////////////////////////////////////////////
import btnStyle from '../styles/button.module.scss';

// ////////////////////////////////////////////////////////////////////////////
// Interfaces and Types
// ////////////////////////////////////////////////////////////////////////////
interface ITeamModalProps {
  readonly team?: ITeam;
  readonly setTeams: React.Dispatch<React.SetStateAction<ITeam[]>>;
  readonly anchor: string | JSX.Element;
}

// ////////////////////////////////////////////////////////////////////////////
// Config
// ////////////////////////////////////////////////////////////////////////////
Modal.setAppElement( document.getElementById( 'root' ) as HTMLElement );

// ////////////////////////////////////////////////////////////////////////////
// Implementation
// ////////////////////////////////////////////////////////////////////////////
export const TeamModal: FC<ITeamModalProps> = ( { team, setTeams, anchor }: ITeamModalProps ) => {
  // Team Setup
  const [localTeam, setLocalTeam] = useState<Partial<ITeam>>( team || { active: true } );

  // Modal Setup
  const [modalIsOpen, setModalIsOpen] = useState( false );

  // Modal Controls
  const openModal = () => setModalIsOpen( true );
  const closeModal = () => setModalIsOpen( false );

  // Update Controls
  const handleUpdate = ( key: string, value: any ) => {
    setLocalTeam( { ...localTeam, [key]: value } );
  };

  // Update Team
  const handleSubmit = async () => {
    const { name, id, active } = localTeam;

    if ( !name ) {
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
  };

  // Styles
  const modalStyle = {
    content: {
      height: 'fit-content',
      width: 'fit-content',
      top: '50%',
      left: '50%',
      right: 'auto',
      bottom: 'auto',
      marginRight: '-50%',
      transform: 'translate(-50%, -50%)',
      boxShadow: 'rgba(0, 0, 0, 0.35) 0px 5px 15px',
    },
  };

  return (
    <>
      <button className={ btnStyle['anchor-btn'] } onClick={ openModal } type="button">{ anchor }</button>
      <Modal
        isOpen={ modalIsOpen }
        onRequestClose={ closeModal }
        contentLabel="Example Modal"
        style={ modalStyle }
      >
        <h1>{ localTeam.id ? `Update ${localTeam.name}` : 'Add a New Team' }</h1>
        <label
          style={ { margin: '0.5rem 0', display: 'block' } }
        >
          Team Name
        </label>
        <input
          style={ { maxWidth: '100%', padding: '0.3rem 0.5rem', display: 'block' } }
          type="text"
          value={ localTeam.name || '' }
          onChange={ e => handleUpdate( 'name', e.target.value ) }
          aria-label="Team Name"
        />
        { localTeam.id && (
          <div style={ { margin: '0.5rem 0', display: 'block' } }>
            <ToggleSwitch
              active={ localTeam.active ?? false }
              callback={ e => handleUpdate( 'active', e ) }
              id={ localTeam.id }
            />
          </div>
        ) }
        <div style={ { margin: '0.5rem 0' } }>
          <button
            className={ `${btnStyle.btn} ${btnStyle['spaced-btn']}` }
            onClick={ handleSubmit }
            type="button"
          >
            Submit
          </button>
          <button
            className={ `${btnStyle.btn} ${btnStyle['spaced-btn']} ${btnStyle['back-btn']} ` }
            onClick={ closeModal }
            type="button"
          >
            Cancel
          </button>
        </div>
      </Modal>
    </>
  );
};
